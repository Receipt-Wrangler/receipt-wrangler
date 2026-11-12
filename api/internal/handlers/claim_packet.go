package handlers

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// GetClaimPacket answers GET /claim/{claimId}/packet with the claim's packet, a
// zip (see services.ClaimPacketService). Gated like reading the claim:
// group.receipts.read in the claim's group, a missing claim the same 403.
//
// The zip is built into temp/ before anything is written, so every failure is
// still a clean JSON error rather than a truncated download: 400 over a size
// bound, 403 for a claim the caller may not see, 500 otherwise.
func GetClaimPacket(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error building claim packet", permissions.GroupReceiptsRead,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			packet, err := services.NewClaimPacketService(nil).BuildClaimPacket(r.Context(), structs.GetClaims(r).UserId, claimId)
			var tooLarge services.ClaimPacketTooLargeError
			if errors.As(err, &tooLarge) {
				utils.WriteCustomErrorResponse(w, tooLarge.Message, http.StatusBadRequest)
				return 0, nil
			}
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			defer packet.Remove()

			file, err := os.Open(packet.Path)
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}
			defer file.Close()

			w.Header().Set("Content-Type", constants.ApplicationZip)
			w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": packet.Filename}))
			w.Header().Set("Content-Length", strconv.FormatInt(packet.Size, 10))
			w.WriteHeader(http.StatusOK)

			// Streaming has begun, so a failure can no longer become an error
			// response; it is logged and the client sees a short download.
			if _, err := io.Copy(w, file); err != nil {
				logging.LogStd(logging.LOG_LEVEL_ERROR, "claim packet: streaming failed: "+err.Error())
			}
			return 0, nil
		})
}
