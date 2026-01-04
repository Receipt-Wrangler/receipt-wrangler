import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { ApiModule, Receipt } from "../../open-api";
import { ReceiptsTableComponent } from "./receipts-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptsTableComponent", () => {
  let component: ReceiptsTableComponent;
  let fixture: ComponentFixture<ReceiptsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ReceiptsTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        NgxsModule.forRoot([ReceiptTableState]),
        ReactiveFormsModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatDialogModule,
        PipesModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        categories: [],
                        tags: [],
                    },
                },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(ReceiptsTableComponent);
    component = fixture.componentInstance;
    component.table = {
      selection: {},
      changed: of(undefined),
    } as any;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should map selected ids from selecton", () => {
    const selectedReceipts: Receipt[] = [
      {
        id: 1,
      } as Receipt,
      {
        id: 2,
      } as Receipt,
    ];
    component.table = {
      selection: {
        changed: of({
          source: {
            selected: selectedReceipts,
          },
        }),
      },
    } as any;
    component.ngAfterViewInit();

    expect(component.selectedReceiptIds).toEqual([1, 2]);
  });
});
