import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { NgxsModule, Store } from "@ngxs/store";
import { PipesModule } from "../../pipes";
import { AuthState } from "../../store";
import { GroupMemberFormComponent } from "./group-member-form.component";

describe("GroupMemberFormComponent", () => {
  let component: GroupMemberFormComponent;
  let fixture: ComponentFixture<GroupMemberFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GroupMemberFormComponent],
      imports: [
        NgxsModule.forRoot([AuthState]),
        PipesModule,
        ReactiveFormsModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        {
          provide: MatDialogRef,
          useValue: {},
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(GroupMemberFormComponent);
    component = fixture.componentInstance;
    TestBed.inject(Store).reset({
      auth: { userId: "1" },
    });
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
