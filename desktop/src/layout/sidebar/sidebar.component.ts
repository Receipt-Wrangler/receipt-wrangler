import { Component, OnInit, ViewEncapsulation } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { Select, Store } from "@ngxs/store";
import { map, Observable, switchMap, take } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { SetPage } from "src/store/receipt-table.actions";
import { AboutComponent } from "../../about/about/about.component";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { ImportFormComponent } from "../../import/import-form/import-form.component";
import { AuthService, Group, GroupStatus, User } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState, GroupState, Logout, SetSelectedGroupId } from "../../store";

@Component({
    selector: "app-sidebar",
    templateUrl: "./sidebar.component.html",
    styleUrls: ["./sidebar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SidebarComponent implements OnInit {
  constructor(
    private authService: AuthService,
    private matDialog: MatDialog,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {}

  @Select(AuthState.loggedInUser) public loggedInUser!: Observable<User>;

  @Select(AuthState.isLoggedIn) public isLoggedIn!: Observable<boolean>;

  @Select(LayoutState.isSidebarOpen) public isSidebarOpen!: Observable<boolean>;

  @Select(GroupState.selectedGroupId)
  public selectedGroupId!: Observable<string>;

  public groups!: Observable<Group[]>;

  public addButtonExpanded: boolean | null = null;


  public ngOnInit(): void {
    this.groups = this.store
      .select(GroupState.groups)
      .pipe(map((g) => g.filter((g) => g.status !== GroupStatus.Archived)));
  }

  public groupClicked(groupId: number): void {
    this.store.dispatch(new SetSelectedGroupId(groupId.toString()));
    this.store.dispatch(new SetPage(1));
    const dashboardLink = this.store.selectSnapshot(GroupState.dashboardLink);
    this.router.navigate([dashboardLink]);
  }

  public toggleAddButtonExpanded(): void {
    this.addButtonExpanded = !this.addButtonExpanded;
  }

  public logout(): void {
    this.authService
      .logout()
      .pipe(
        take(1),
        switchMap(() => this.store.dispatch(new Logout())),
        switchMap(() => this.router.navigate(["/"])),
      )
      .subscribe();
  }

  public openImportDialog(): void {
    this.matDialog.open(ImportFormComponent, DEFAULT_DIALOG_CONFIG);
  }

  public openAboutDialog(): void {
    this.matDialog.open(AboutComponent, DEFAULT_DIALOG_CONFIG);
  }
}
