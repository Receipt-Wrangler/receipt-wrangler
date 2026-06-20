import { Component, OnInit, Signal, input, output, signal } from "@angular/core";
import { FormArray, FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { Comment, CommentService, Permission } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState } from "../../store";

@UntilDestroy()
@Component({
    selector: "app-receipt-comments",
    templateUrl: "./receipt-comments.component.html",
    styleUrls: ["./receipt-comments.component.scss"],
    standalone: false
})
export class ReceiptCommentsComponent implements OnInit {
  loggedInUserId = this.store.selectSignal(AuthState.userId);
  public readonly comments = input<Comment[]>([]);
  public internalComments = signal<Comment[]>([]);
  public readonly mode = input.required<FormMode>();
  public readonly receiptId = input<number>();
  public readonly groupId = input<number>();
  public readonly commentsUpdated = output<FormArray>();

  public canCreateComments!: Signal<boolean>;

  public canDeleteComments!: Signal<boolean>;

  public formMode = FormMode;

  public commentsArray: FormArray<FormGroup> = new FormArray<FormGroup>([]);

  public newCommentFormControl: FormControl = new FormControl("");


  constructor(
    private formBuilder: FormBuilder,
    private store: Store,
    private commentService: CommentService,
    private snackbarService: SnackbarService
  ) {}

  public ngOnInit(): void {
    this.internalComments.set(this.comments());
    this.setCommentPermissions();
    this.initForm();
  }

  private setCommentPermissions(): void {
    const groupId = this.groupId() ?? 0;
    this.canCreateComments = this.store.selectSignal(
      AuthState.hasGroupPermission(groupId, Permission.GroupCommentsCreate)
    );
    this.canDeleteComments = this.store.selectSignal(
      AuthState.hasGroupPermission(groupId, Permission.GroupCommentsDelete)
    );
  }

  private initForm(): void {
    this.internalComments().forEach((c) => {
      this.commentsArray.push(this.buildCommentFormGroup(c));
    });
  }

  private buildCommentFormGroup(comment?: Comment): FormGroup {
    return this.formBuilder.group({
      comment: [comment?.comment ?? "", Validators.required],
      userId: [
        comment?.userId ??
        Number.parseInt(this.store.selectSnapshot(AuthState.userId)),
      ],
      receiptId: [comment?.receiptId ?? this.receiptId()],
    });
  }

  public addComment(): void {
    const isValid = this.newCommentFormControl.valid && this.newCommentFormControl.value.trim() !== "";
    const newComment = {
      comment: this.newCommentFormControl.value,
      userId: Number.parseInt(this.store.selectSnapshot(AuthState.userId)),
      receiptId: this.receiptId(),
    } as any;

    const mode = this.mode();
    if (isValid && mode === FormMode.add) {
      this.commentsArray.push(this.buildCommentFormGroup(newComment));
      this.newCommentFormControl.reset();
      this.commentsUpdated.emit(this.commentsArray);
    } else if (isValid && mode === FormMode.edit) {
      this.commentService
        .addComment(newComment)
        .pipe(
          take(1),
          tap((comment: Comment) => {
            this.internalComments.update(comments => [...comments, comment]);
            this.commentsArray.push(this.buildCommentFormGroup(newComment));
            this.snackbarService.success("Comment successfully added");
            this.newCommentFormControl.reset();
          })
        )
        .subscribe();
    }
  }

  public deleteComment(index: number): void {
    switch (this.mode()) {
      case FormMode.edit:
      case FormMode.view:
        const comment = this.internalComments()[index];
        let commentIdToDelete = comment.id;

        this.commentService
          .deleteComment(commentIdToDelete)
          .pipe(
            take(1),
            tap(() => {
              this.commentsArray.removeAt(index);
              this.internalComments.set(this.internalComments().filter(
                (c) => c.id !== comment.id
              ));
              this.snackbarService.success("Comment successfully deleted");
            })
          )
          .subscribe();
        break;

      case FormMode.add:
        this.commentsArray.removeAt(index);
        break;
    }
  }
}
