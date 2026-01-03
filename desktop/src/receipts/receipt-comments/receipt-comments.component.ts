import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core";
import { FormArray, FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { Observable, take, tap } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { Comment, CommentService } from "../../open-api";
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
  @Select(AuthState.userId) public loggedInUserId!: Observable<string>;
  @Input() public comments: Comment[] = [];
  @Input() public mode!: FormMode;
  @Input() public receiptId?: number;
  @Output() public commentsUpdated: EventEmitter<FormArray> =
    new EventEmitter<FormArray>();

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
    this.initForm();
  }

  private initForm(): void {
    this.comments.forEach((c) => {
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
      receiptId: [comment?.receiptId ?? this.receiptId],
    });
  }

  public addComment(): void {
    const isValid = this.newCommentFormControl.valid && this.newCommentFormControl.value.trim() !== "";
    const newComment = {
      comment: this.newCommentFormControl.value,
      userId: Number.parseInt(this.store.selectSnapshot(AuthState.userId)),
      receiptId: this.receiptId,
    } as any;

    if (isValid && this.mode === FormMode.add) {
      this.commentsArray.push(this.buildCommentFormGroup(newComment));
      this.newCommentFormControl.reset();
      this.commentsUpdated.emit(this.commentsArray);
    } else if (isValid && this.mode === FormMode.edit) {
      this.commentService
        .addComment(newComment)
        .pipe(
          take(1),
          tap((comment: Comment) => {
            this.comments.push(comment);
            this.commentsArray.push(this.buildCommentFormGroup(newComment));
            this.snackbarService.success("Comment successfully added");
            this.newCommentFormControl.reset();
          })
        )
        .subscribe();
    }
  }

  public deleteComment(index: number): void {
    switch (this.mode) {
      case FormMode.edit:
      case FormMode.view:
        const comment = this.comments[index];
        let commentIdToDelete = comment.id;

        this.commentService
          .deleteComment(commentIdToDelete)
          .pipe(
            take(1),
            tap(() => {
              this.commentsArray.removeAt(index);
              this.comments = this.comments.filter(
                (c) => c.id !== comment.id
              );
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
