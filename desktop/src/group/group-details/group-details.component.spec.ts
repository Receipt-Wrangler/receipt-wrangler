import { CommonModule } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { AuthState, GroupState } from "../../store/index";

import { GroupDetailsComponent } from "./group-details.component";

describe("GroupDetailsComponent", () => {
  let component: GroupDetailsComponent;
  let fixture: ComponentFixture<GroupDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GroupDetailsComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ReactiveFormsModule, CommonModule, NgxsModule.forRoot([GroupState, AuthState])],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                group: { id: 1 },
              },
            },
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(GroupDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
