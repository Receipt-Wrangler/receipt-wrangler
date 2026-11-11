import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { LayoutState } from "src/store/layout.state";
import { SubmitButtonComponent } from "./submit-button.component";

describe("SubmitButtonComponent", () => {
  let component: SubmitButtonComponent;
  let fixture: ComponentFixture<SubmitButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([LayoutState]), SubmitButtonComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {},
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(SubmitButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
