package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// This file is the receipt history's diff: pure functions from two receipt
// snapshots to the change rows a client formats. receipt_history.go loads,
// selects and redacts the snapshots; nothing here touches the database.

// receiptDiffContext supplies what a diff shows but a snapshot may not carry.
type receiptDiffContext struct {
	// groupNames names every group a snapshot sits in, already redacted: a group
	// the caller may not see is "(Restricted)".
	groupNames map[uint]string
	// customFields are the stored definitions, by id. A snapshot's own
	// definition wins, since it holds the name and option text of the time;
	// these cover the values a version 1 "before" stored without one.
	customFields map[uint]models.CustomField
}

// diffReceiptSnapshots lists what changed from before to after, in a fixed
// field order. A nil before is the creation: after is compared with an empty
// receipt, so only the fields it set appear, each ADDED, with the email origin
// first. BaseModel fields (id, timestamps, created-by) are never compared.
func diffReceiptSnapshots(before *models.Receipt, after models.Receipt, context receiptDiffContext) []structs.ReceiptHistoryChange {
	diff := receiptDiff{
		context:  context,
		creation: before == nil,
		changes:  make([]structs.ReceiptHistoryChange, 0),
	}
	from := models.Receipt{}
	if before != nil {
		from = *before
	}

	if diff.creation {
		diff.compare(structs.ReceiptHistoryFieldEmailInbox, "", textSide(""), textSide(after.EmailInbox))
		diff.compare(structs.ReceiptHistoryFieldEmailSender, "", textSide(""), textSide(after.EmailSender))
	}

	diff.compare(structs.ReceiptHistoryFieldName, "", textSide(from.Name), textSide(after.Name))
	diff.compare(structs.ReceiptHistoryFieldAmount, "", currencySide(&from.Amount), currencySide(&after.Amount))
	diff.compare(structs.ReceiptHistoryFieldDate, "", calendarDaySide(&from.Date), calendarDaySide(&after.Date))
	diff.compare(structs.ReceiptHistoryFieldGroup, "", diff.groupSide(from.GroupId), diff.groupSide(after.GroupId))
	diff.compare(structs.ReceiptHistoryFieldPaidBy, "", userSide(from.PaidByUserID), userSide(after.PaidByUserID))
	diff.compare(structs.ReceiptHistoryFieldStatus, "", statusSide(string(from.Status)), statusSide(string(after.Status)))
	diff.compare(structs.ReceiptHistoryFieldResolvedDate, "", instantSide(from.ResolvedDate), instantSide(after.ResolvedDate))
	diff.append(namedSetChanges(structs.ReceiptHistoryFieldCategories, "", categoryRefs(from.Categories), categoryRefs(after.Categories)))
	diff.append(namedSetChanges(structs.ReceiptHistoryFieldTags, "", tagRefs(from.Tags), tagRefs(after.Tags)))
	diff.customFields(from.CustomFields, after.CustomFields)
	diff.items(from.ReceiptItems, after.ReceiptItems)
	diff.comments(from.Comments, after.Comments)
	diff.images(from.ImageFiles, after.ImageFiles)

	return diff.changes
}

// receiptHistoryAction names what an update entry did, from its changes (read
// after redaction, so it describes what the caller is shown).
//
// Only status-type changes — the receipt's status, its resolved date, its
// shares' statuses (which a receipt status cascades to) and a comment (which a
// bulk status change may add) — with the status among them is a status change,
// RESOLVED when the receipt became resolved. Only comments and images is
// IMAGES_COMMENTS_UPDATED. Anything else, including no visible change, is
// UPDATED.
func receiptHistoryAction(changes []structs.ReceiptHistoryChange) string {
	if len(changes) == 0 {
		return structs.ReceiptHistoryActionUpdated
	}

	statusOnly := true
	imagesCommentsOnly := true
	var newStatus *structs.ReceiptHistoryValue
	for _, change := range changes {
		switch {
		case change.Field == structs.ReceiptHistoryFieldStatus:
			newStatus = change.To
			imagesCommentsOnly = false
		case change.Field == structs.ReceiptHistoryFieldResolvedDate, isShareStatusChange(change):
			imagesCommentsOnly = false
		case change.Field == structs.ReceiptHistoryFieldComment:
		case change.Field == structs.ReceiptHistoryFieldImage:
			statusOnly = false
		default:
			statusOnly = false
			imagesCommentsOnly = false
		}
	}

	switch {
	case statusOnly && newStatus != nil && newStatus.Value == string(models.RESOLVED):
		return structs.ReceiptHistoryActionResolved
	case statusOnly && newStatus != nil:
		return structs.ReceiptHistoryActionStatusChanged
	case imagesCommentsOnly:
		return structs.ReceiptHistoryActionImagesCommentsUpdated
	default:
		return structs.ReceiptHistoryActionUpdated
	}
}

func isShareStatusChange(change structs.ReceiptHistoryChange) bool {
	return change.Field == structs.ReceiptHistoryFieldShare &&
		change.To != nil && change.To.Type == structs.ReceiptHistoryValueStatus
}

type receiptDiff struct {
	context  receiptDiffContext
	creation bool
	changes  []structs.ReceiptHistoryChange
}

// diffSide is one side of a comparison: key is what is compared, value what is
// shown (nil when the side has no value). Keeping them apart lets a select be
// compared by option id while showing its text, and a group by id while
// showing its name.
type diffSide struct {
	key   string
	value *structs.ReceiptHistoryValue
}

// compare records a change when the two sides' keys differ: ADDED when the
// "from" side has no value (always, on the creation), REMOVED when the "to"
// side has none, CHANGED otherwise.
func (diff *receiptDiff) compare(field string, label string, from diffSide, to diffSide) {
	diff.compareWithUser(field, label, nil, from, to)
}

func (diff *receiptDiff) compareWithUser(field string, label string, userId *uint, from diffSide, to diffSide) {
	if from.key == to.key {
		return
	}
	if diff.creation {
		from.value = nil
	}
	if from.value == nil && to.value == nil {
		return
	}

	diff.changes = append(diff.changes, newChange(field, label, userId, from.value, to.value))
}

func (diff *receiptDiff) append(changes []structs.ReceiptHistoryChange) {
	diff.changes = append(diff.changes, changes...)
}

func newChange(field string, label string, userId *uint, from *structs.ReceiptHistoryValue, to *structs.ReceiptHistoryValue) structs.ReceiptHistoryChange {
	kind := structs.ReceiptHistoryKindChanged
	switch {
	case from == nil:
		kind = structs.ReceiptHistoryKindAdded
	case to == nil:
		kind = structs.ReceiptHistoryKindRemoved
	}

	return structs.ReceiptHistoryChange{
		Field:  field,
		Kind:   kind,
		Label:  label,
		UserId: copyUintPointer(userId),
		From:   from,
		To:     to,
	}
}

func (diff *receiptDiff) groupSide(groupId uint) diffSide {
	if groupId == 0 {
		return diffSide{}
	}
	name, ok := diff.context.groupNames[groupId]
	if !ok {
		name = restrictedGroupName
	}
	return diffSide{key: strconv.FormatUint(uint64(groupId), 10), value: textValue(name)}
}

// customFields compares values by custom field id. An absent value equals an
// empty one, so a default the form attached empty makes no change.
func (diff *receiptDiff) customFields(from []models.CustomFieldValue, to []models.CustomFieldValue) {
	fromById := customFieldValuesById(from)
	toById := customFieldValuesById(to)

	ids := make([]uint, 0, len(fromById)+len(toById))
	for id := range fromById {
		ids = append(ids, id)
	}
	for id := range toById {
		if _, ok := fromById[id]; !ok {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		fromValue, toValue := fromById[id], toById[id]
		fromDefinition, fromOk := diff.customFieldDefinition(id, fromValue, toValue)
		toDefinition, toOk := diff.customFieldDefinition(id, toValue, fromValue)
		if !fromOk || !toOk {
			continue
		}

		diff.compare(
			structs.ReceiptHistoryFieldCustomField,
			toDefinition.Name,
			diff.customFieldSide(fromDefinition, fromValue),
			diff.customFieldSide(toDefinition, toValue),
		)
	}
}

// customFieldDefinition is the definition a side's value is read against: its
// own snapshot's, else the other side's, else the stored one.
func (diff *receiptDiff) customFieldDefinition(id uint, own *models.CustomFieldValue, other *models.CustomFieldValue) (models.CustomField, bool) {
	for _, value := range []*models.CustomFieldValue{own, other} {
		if value != nil && value.CustomField.Type != "" {
			return value.CustomField, true
		}
	}
	definition, ok := diff.context.customFields[id]
	return definition, ok && definition.Type != ""
}

func customFieldValuesById(values []models.CustomFieldValue) map[uint]*models.CustomFieldValue {
	byId := make(map[uint]*models.CustomFieldValue, len(values))
	for i := range values {
		byId[values[i].CustomFieldId] = &values[i]
	}
	return byId
}

func (diff *receiptDiff) customFieldSide(definition models.CustomField, value *models.CustomFieldValue) diffSide {
	if value == nil {
		value = &models.CustomFieldValue{}
	}

	switch definition.Type {
	case models.TEXT:
		if value.StringValue == nil {
			return diffSide{}
		}
		return textSide(*value.StringValue)
	case models.DATE:
		return calendarDaySide(value.DateValue)
	case models.SELECT:
		if value.SelectValue == nil || *value.SelectValue == 0 {
			return diffSide{}
		}
		return diffSide{
			key:   strconv.FormatUint(uint64(*value.SelectValue), 10),
			value: textValue(diff.customFieldOptionText(definition, *value.SelectValue)),
		}
	case models.CURRENCY:
		return currencySide(value.CurrencyValue)
	case models.BOOLEAN:
		// An unset checkbox and an unticked one read the same, so they compare
		// equal; only a set value is shown.
		checked := value.BooleanValue != nil && *value.BooleanValue
		side := diffSide{key: strconv.FormatBool(checked)}
		if value.BooleanValue != nil {
			side.value = &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueBoolean, Value: side.key}
		}
		return side
	default:
		return diffSide{}
	}
}

// customFieldOptionText is a select value's option text: the definition's own
// (the text of the time), else the stored one.
func (diff *receiptDiff) customFieldOptionText(definition models.CustomField, optionId uint) string {
	for _, options := range [][]models.CustomFieldOption{definition.Options, diff.context.customFields[definition.ID].Options} {
		for _, option := range options {
			if option.ID == optionId {
				return option.Value
			}
		}
	}
	return ""
}

// itemKey is what an item is matched on. Every save deletes and recreates a
// receipt's items, so their ids never match across two snapshots.
type itemKey struct {
	name      string
	chargedTo uint
}

type itemPair struct {
	from    models.Item
	to      models.Item
	toIndex int
}

// items compares items (charged to no one) and shares (charged to a member)
// as a multiset on (trimmed name, charged-to user). A pair reports a changed
// amount and changed categories/tags; a share also reports a changed status.
// A plain item's status is not reported: the receipt's status cascades to
// every item, and the receipt's own status change already says so.
func (diff *receiptDiff) items(from []models.Item, to []models.Item) {
	fromByKey := itemsByKey(flattenReceiptItems(from))
	toByKey := itemsByKey(flattenReceiptItems(to))

	keys := make([]itemKey, 0, len(fromByKey)+len(toByKey))
	for key := range fromByKey {
		keys = append(keys, key)
	}
	for key := range toByKey {
		if _, ok := fromByKey[key]; !ok {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].name != keys[j].name {
			return keys[i].name < keys[j].name
		}
		return keys[i].chargedTo < keys[j].chargedTo
	})

	var itemChanges, shareChanges, categoryChanges, tagChanges []structs.ReceiptHistoryChange
	for _, key := range keys {
		field := structs.ReceiptHistoryFieldItem
		var chargedTo *uint
		if key.chargedTo != 0 {
			field = structs.ReceiptHistoryFieldShare
			chargedTo = &key.chargedTo
		}

		pairs, added, removed := matchItems(fromByKey[key], toByKey[key])
		changes := make([]structs.ReceiptHistoryChange, 0)
		for _, pair := range pairs {
			if !pair.from.Amount.Equal(pair.to.Amount) {
				changes = append(changes, newChange(field, key.name, chargedTo, currencyValue(pair.from.Amount), currencyValue(pair.to.Amount)))
			}
			if chargedTo != nil && pair.from.Status != pair.to.Status {
				changes = append(changes, newChange(field, key.name, chargedTo, statusSide(string(pair.from.Status)).value, statusSide(string(pair.to.Status)).value))
			}
			categoryChanges = append(categoryChanges, namedSetChanges(structs.ReceiptHistoryFieldItemCategories, key.name, categoryRefs(pair.from.Categories), categoryRefs(pair.to.Categories))...)
			tagChanges = append(tagChanges, namedSetChanges(structs.ReceiptHistoryFieldItemTags, key.name, tagRefs(pair.from.Tags), tagRefs(pair.to.Tags))...)
		}
		for _, item := range added {
			changes = append(changes, newChange(field, key.name, chargedTo, nil, currencyValue(item.Amount)))
		}
		for _, item := range removed {
			changes = append(changes, newChange(field, key.name, chargedTo, currencyValue(item.Amount), nil))
		}

		if chargedTo == nil {
			itemChanges = append(itemChanges, changes...)
		} else {
			shareChanges = append(shareChanges, changes...)
		}
	}

	diff.append(itemChanges)
	diff.append(shareChanges)
	diff.append(categoryChanges)
	diff.append(tagChanges)
}

// flattenReceiptItems lists a receipt's items with their linked items, once
// each: a version 1 snapshot lists linked items at the top level too.
func flattenReceiptItems(items []models.Item) []models.Item {
	flattened := make([]models.Item, 0, len(items))
	seen := make(map[uint]bool, len(items))
	add := func(item models.Item) {
		if item.ID != 0 {
			if seen[item.ID] {
				return
			}
			seen[item.ID] = true
		}
		flattened = append(flattened, item)
	}

	for _, item := range items {
		add(item)
		for _, linkedItem := range item.LinkedItems {
			add(linkedItem)
		}
	}
	return flattened
}

func itemsByKey(items []models.Item) map[itemKey][]models.Item {
	byKey := make(map[itemKey][]models.Item)
	for _, item := range items {
		key := itemKey{name: strings.TrimSpace(item.Name)}
		if item.ChargedToUserId != nil {
			key.chargedTo = *item.ChargedToUserId
		}
		byKey[key] = append(byKey[key], item)
	}
	return byKey
}

// matchItems pairs the items sharing one key: equal amounts first, then the
// rest in order, which differ in amount. What is left over was added or
// removed. Pairs come back in the "to" side's order.
func matchItems(from []models.Item, to []models.Item) ([]itemPair, []models.Item, []models.Item) {
	fromUsed := make([]bool, len(from))
	toUsed := make([]bool, len(to))
	pairs := make([]itemPair, 0)

	for j := range to {
		for i := range from {
			if !fromUsed[i] && from[i].Amount.Equal(to[j].Amount) {
				pairs = append(pairs, itemPair{from: from[i], to: to[j], toIndex: j})
				fromUsed[i], toUsed[j] = true, true
				break
			}
		}
	}

	i := 0
	for j := range to {
		if toUsed[j] {
			continue
		}
		for i < len(from) && fromUsed[i] {
			i++
		}
		if i == len(from) {
			break
		}
		pairs = append(pairs, itemPair{from: from[i], to: to[j], toIndex: j})
		fromUsed[i], toUsed[j] = true, true
	}
	sort.SliceStable(pairs, func(a, b int) bool { return pairs[a].toIndex < pairs[b].toIndex })

	added := make([]models.Item, 0)
	for j, used := range toUsed {
		if !used {
			added = append(added, to[j])
		}
	}
	removed := make([]models.Item, 0)
	for i, used := range fromUsed {
		if !used {
			removed = append(removed, from[i])
		}
	}

	return pairs, added, removed
}

// comments compares comments, replies included, by id.
func (diff *receiptDiff) comments(from []models.Comment, to []models.Comment) {
	fromById := commentsById(from)
	toById := commentsById(to)

	for _, id := range unionSortedIds(fromById, toById) {
		fromComment, inFrom := fromById[id]
		toComment, inTo := toById[id]

		var fromSide, toSide diffSide
		author := (*uint)(nil)
		if inFrom {
			fromSide = diffSide{key: "present:" + fromComment.Comment, value: textValue(fromComment.Comment)}
			author = fromComment.UserId
		}
		if inTo {
			toSide = diffSide{key: "present:" + toComment.Comment, value: textValue(toComment.Comment)}
			author = toComment.UserId
		}
		diff.compareWithUser(structs.ReceiptHistoryFieldComment, "", author, fromSide, toSide)
	}
}

func commentsById(comments []models.Comment) map[uint]models.Comment {
	byId := make(map[uint]models.Comment)
	var visit func(comments []models.Comment)
	visit = func(comments []models.Comment) {
		for _, comment := range comments {
			if _, seen := byId[comment.ID]; !seen {
				byId[comment.ID] = comment
			}
			visit(comment.Replies)
		}
	}
	visit(comments)
	return byId
}

// images compares image files by id.
func (diff *receiptDiff) images(from []models.FileData, to []models.FileData) {
	fromById := make(map[uint]models.FileData, len(from))
	for _, file := range from {
		fromById[file.ID] = file
	}
	toById := make(map[uint]models.FileData, len(to))
	for _, file := range to {
		toById[file.ID] = file
	}

	for _, id := range unionSortedIds(fromById, toById) {
		var fromSide, toSide diffSide
		if file, ok := fromById[id]; ok {
			fromSide = diffSide{key: "present", value: textValue(file.Name)}
		}
		if file, ok := toById[id]; ok {
			toSide = diffSide{key: "present", value: textValue(file.Name)}
		}
		diff.compare(structs.ReceiptHistoryFieldImage, "", fromSide, toSide)
	}
}

func unionSortedIds[T any](a map[uint]T, b map[uint]T) []uint {
	ids := make([]uint, 0, len(a)+len(b))
	for id := range a {
		ids = append(ids, id)
	}
	for id := range b {
		if _, ok := a[id]; !ok {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

// namedRef is a category or tag: compared by id, shown by name.
type namedRef struct {
	id   uint
	name string
}

func categoryRefs(categories []models.Category) []namedRef {
	refs := make([]namedRef, len(categories))
	for i, category := range categories {
		refs[i] = namedRef{id: category.ID, name: category.Name}
	}
	return refs
}

func tagRefs(tags []models.Tag) []namedRef {
	refs := make([]namedRef, len(tags))
	for i, tag := range tags {
		refs[i] = namedRef{id: tag.ID, name: tag.Name}
	}
	return refs
}

// namedSetChanges lists one ADDED row per id only in "to", then one REMOVED
// row per id only in "from", each by name.
func namedSetChanges(field string, label string, from []namedRef, to []namedRef) []structs.ReceiptHistoryChange {
	changes := make([]structs.ReceiptHistoryChange, 0)
	for _, ref := range namedRefsMissingFrom(to, from) {
		changes = append(changes, newChange(field, label, nil, nil, textValue(ref.name)))
	}
	for _, ref := range namedRefsMissingFrom(from, to) {
		changes = append(changes, newChange(field, label, nil, textValue(ref.name), nil))
	}
	return changes
}

// namedRefsMissingFrom returns the refs whose id other lacks, by name.
func namedRefsMissingFrom(refs []namedRef, other []namedRef) []namedRef {
	otherIds := make(map[uint]bool, len(other))
	for _, ref := range other {
		otherIds[ref.id] = true
	}

	missing := make([]namedRef, 0)
	seen := make(map[uint]bool, len(refs))
	for _, ref := range refs {
		if otherIds[ref.id] || seen[ref.id] {
			continue
		}
		seen[ref.id] = true
		missing = append(missing, ref)
	}
	sort.Slice(missing, func(i, j int) bool {
		if missing[i].name != missing[j].name {
			return missing[i].name < missing[j].name
		}
		return missing[i].id < missing[j].id
	})
	return missing
}

func textValue(text string) *structs.ReceiptHistoryValue {
	return &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueText, Value: text}
}

func currencyValue(amount decimal.Decimal) *structs.ReceiptHistoryValue {
	return &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueCurrency, Value: amount.String()}
}

func textSide(text string) diffSide {
	if text == "" {
		return diffSide{}
	}
	return diffSide{key: text, value: textValue(text)}
}

// currencySide compares amounts by value, so 12.5 and 12.50 are equal.
func currencySide(amount *decimal.Decimal) diffSide {
	if amount == nil {
		return diffSide{}
	}
	value := currencyValue(*amount)
	return diffSide{key: value.Value, value: value}
}

// calendarDaySide compares a calendar day stored as midnight UTC by its UTC
// day, which is how it is written and read everywhere.
func calendarDaySide(day *time.Time) diffSide {
	if day == nil || day.IsZero() {
		return diffSide{}
	}
	formatted := day.UTC().Format(time.DateOnly)
	return diffSide{key: formatted, value: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueDate, Value: formatted}}
}

func instantSide(instant *time.Time) diffSide {
	if instant == nil || instant.IsZero() {
		return diffSide{}
	}
	formatted := instant.UTC().Format(time.RFC3339)
	return diffSide{key: formatted, value: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueInstant, Value: formatted}}
}

func userSide(userId uint) diffSide {
	if userId == 0 {
		return diffSide{}
	}
	id := userId
	return diffSide{
		key:   strconv.FormatUint(uint64(userId), 10),
		value: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueUser, UserId: &id},
	}
}

func statusSide(status string) diffSide {
	if status == "" {
		return diffSide{}
	}
	return diffSide{key: status, value: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueStatus, Value: status}}
}

func copyUintPointer(value *uint) *uint {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
