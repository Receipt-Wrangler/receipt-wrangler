import { CurrencyPipe } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatCardModule } from "@angular/material/card";
import { MatListModule } from "@angular/material/list";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { ApiModule, UserService } from "../../open-api";
import { UserState } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { CardComponent } from "../card/card.component";
import { SummaryCardComponent } from "./summary-card.component";

describe("SummaryCardComponent", () => {
  let component: SummaryCardComponent;
  let fixture: ComponentFixture<SummaryCardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ApiModule,
        MatCardModule,
        MatListModule,
        NgxsModule.forRoot([UserState, SystemSettingsState]),
        SummaryCardComponent, CardComponent],
      providers: [
        CurrencyPipe,
        provideZonelessChangeDetection(),
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({}),
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(SummaryCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should set user data correctly when there is data", async () => {
    const usersService = TestBed.inject(UserService);
    jest.spyOn(usersService, "getAmountOwedForUser").mockReturnValue(
      of({
        "1": 200,
        "2": -500,
      } as any)
    );

    fixture.componentRef.setInput("groupId", "1");
    await fixture.whenStable();

    expect(Array.from(component.userOwesMap().entries())).toEqual([["1", "200"]]);
    expect(Array.from(component.usersOweMap().entries())).toEqual([["2", "500"]]);
  });
});
