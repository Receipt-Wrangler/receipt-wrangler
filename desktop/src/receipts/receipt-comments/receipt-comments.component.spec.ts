import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { ApiModule, Comment, CommentService } from "../../open-api";
import { AuthState } from "../../store";
import { ReceiptCommentsComponent } from "./receipt-comments.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptCommentsComponent", () => {
  let component: ReceiptCommentsComponent;
  let fixture: ComponentFixture<ReceiptCommentsComponent>;
  let comments: Comment[];
  let store: Store;

  beforeEach(async () => {
    comments = [
      {
        id: 1,
        comment: "comment",
        receiptId: 1,
        userId: 1,
        updatedAt: "",
        createdAt: "",
      },
      {
        id: 2,
        comment: "new comment",
        receiptId: 1,
        userId: 1,
        updatedAt: "",
        createdAt: "",
      },
    ];

    await TestBed.configureTestingModule({
    declarations: [ReceiptCommentsComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        ReactiveFormsModule,
        NgxsModule.forRoot([AuthState]),
        MatSnackBarModule,
        PipesModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
}).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(ReceiptCommentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init each comment correctly", () => {
    component.comments = comments;

    component.ngOnInit();
    expect(component.commentsArray.value).toEqual([
      {
        comment: "comment",
        userId: 1,
        receiptId: 1,
      },
      {
        comment: "new comment",
        userId: 1,
        receiptId: 1,
      },
    ]);
  });


  it("should delete comment that is a top level comment", () => {
    const spy = jest.spyOn(TestBed.inject(CommentService), "deleteComment");
    spy.mockReturnValue(of(undefined as any));
    component.comments = comments;

    component.ngOnInit();
    component.mode = FormMode.view;
    component.deleteComment(0);

    expect(spy).toHaveBeenCalledWith(1);

    expect(component.commentsArray.value.find((c) => c.id === 1)).toEqual(
      undefined
    );
    expect(component.commentsArray.value.length).toEqual(1);
    expect(component.comments.find((c) => c.id === 1)).toEqual(undefined);
    expect(component.comments.length).toEqual(1);
  });


  it("should delete comment in add mode", () => {
    component.ngOnInit();
    component.mode = FormMode.add;
    component.commentsArray.push(new FormGroup({}));

    expect(component.commentsArray.length).toEqual(1);
    expect(component.comments.length).toEqual(0);

    component.deleteComment(0);

    expect(component.commentsArray.length).toEqual(0);
    expect(component.comments.length).toEqual(0);
  });

  it("should add comment if form is valid and is in add mode", () => {
    const eventEmitterSpy = jest.spyOn(component.commentsUpdated, "emit");
    store.reset({
      auth: {
        userId: 1,
      },
    });

    component.mode = FormMode.add;
    component.newCommentFormControl.patchValue("new comment");
    component.receiptId = 1;
    component.addComment();

    expect(component.newCommentFormControl.value).toEqual(null);
    expect(component.newCommentFormControl.pristine).toEqual(true);
    expect(component.commentsArray.length).toEqual(1);
    expect(eventEmitterSpy).toHaveBeenCalledWith(component.commentsArray);
    expect(component.commentsArray.at(0).value).toEqual({
      userId: 1,
      receiptId: 1,
      comment: "new comment",
    });
  });

  it("should send api call if form is valid and is in edit mode", () => {
    const spy = jest.spyOn(TestBed.inject(CommentService), "addComment");
    spy.mockReturnValue(
      of({
        id: 5,
        userId: 1,
        receiptId: 1,
        comment: "new comment",
        createdAt: "",
        updatedAt: "",
      }) as any
    );

    store.reset({
      auth: {
        userId: 1,
      },
    });

    component.mode = FormMode.edit;
    component.newCommentFormControl.patchValue("new comment");
    component.receiptId = 1;
    component.addComment();

    expect(component.commentsArray.length).toEqual(1);
    expect(component.commentsArray.at(0).value).toEqual({
      userId: 1,
      receiptId: 1,
      comment: "new comment",
    });
    expect(component.comments.length).toEqual(1);
    expect(component.comments[0]).toEqual({
      id: 5,
      userId: 1,
      receiptId: 1,
      comment: "new comment",
      createdAt: "",
      updatedAt: "",
    } as any);
  });
});
