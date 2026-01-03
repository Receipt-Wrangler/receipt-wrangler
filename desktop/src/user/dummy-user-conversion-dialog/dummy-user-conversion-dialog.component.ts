import { Component, Input, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { switchMap, take, tap } from "rxjs";
import { User, UserService } from "../../open-api";
import { SnackbarService } from "../../services";
import { UpdateUser } from "../../store";

@Component({
    selector: "app-dummy-user-conversion-dialog",
    templateUrl: "./dummy-user-conversion-dialog.component.html",
    styleUrls: ["./dummy-user-conversion-dialog.component.scss"],
    standalone: false
})
export class DummyUserConversionDialogComponent implements OnInit {
  @Input() public user!: User;

  public form: FormGroup = new FormGroup({});

  constructor(
    private formBuilder: FormBuilder,
    private snackbarService: SnackbarService,
    private store: Store,
    private userService: UserService,
    public matDialogRef: MatDialogRef<DummyUserConversionDialogComponent>
  ) {}

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      password: ["", Validators.required],
    });
  }

  public submitButtonClicked(): void {
    if (this.form.valid) {
      let userId: string = this.user?.id?.toString();
      this.userService
        .convertDummyUserById( Number.parseInt(userId), this.form.value)
        .pipe(
          take(1),
          tap(() => {
            this.snackbarService.success(
              `${this.user.displayName} sucessfully converted to normal user`
            );
          }),
          switchMap(() =>
            this.store.dispatch(
              new UpdateUser(userId, { ...this.user, isDummyUser: false })
            )
          ),
          tap(() => this.matDialogRef.close(true))
        )
        .subscribe();
    }
  }

  public cancelButtonClicked(): void {
    this.matDialogRef.close(false);
  }
}
