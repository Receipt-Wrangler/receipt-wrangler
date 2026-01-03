import { provideHttpClientTesting } from "@angular/common/http/testing";
import {ComponentFixture, TestBed} from "@angular/core/testing";
import {MatTableDataSource} from "@angular/material/table";
import {NoopAnimationsModule} from "@angular/platform-browser/animations";
import {ActivatedRoute} from "@angular/router";
import {NgxsModule} from "@ngxs/store";
import {of} from "rxjs";
import {ConfirmationDialogComponent} from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import {SharedUiModule} from "../../shared-ui/shared-ui.module";
import {SystemEmailTableState} from "../../store/system-email-table.state";
import {TableModule} from "../../table/table.module";

import {SystemEmailTableComponent} from "./system-email-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("SystemEmailsComponent", () => {
  let component: SystemEmailTableComponent;
  let fixture: ComponentFixture<SystemEmailTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [SystemEmailTableComponent],
    imports: [SharedUiModule,
        NgxsModule.forRoot([SystemEmailTableState]),
        TableModule,
        NoopAnimationsModule],
    providers: [{
            provide: ActivatedRoute,
            useValue: {}
        }, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
})
      .compileComponents();

    fixture = TestBed.createComponent(SystemEmailTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should pop confirmation dialog", () => {
    const matDialogSpy = jest.spyOn(component.matDialog, "open").mockReturnValue({
      componentInstance: {},
      afterClosed: () => of(undefined)
    } as any);
    component.dataSource = new MatTableDataSource([
      {
        id: 1,
        username: "test",
      }
    ] as any[]);
    component.deleteButtonClicked(component.dataSource.data[0]);

    expect(matDialogSpy).toHaveBeenCalledWith(ConfirmationDialogComponent);
  });
});
