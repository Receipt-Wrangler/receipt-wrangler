import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { NgxsModule, Store } from "@ngxs/store";
import { FeatureConfigState } from "../../store";
import { GroupTabsComponent } from "./group-tabs.component";

describe("GroupTabsComponent", () => {
  let component: GroupTabsComponent;
  let fixture: ComponentFixture<GroupTabsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [GroupTabsComponent],
      imports: [ReactiveFormsModule, NgxsModule.forRoot([FeatureConfigState])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(GroupTabsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init tabs without group settings", () => {
    component.ngOnInit();

    expect(component.tabs).toEqual([
      {
        label: "Group Details",
        routerLink: "details/view",
        name: "details"
      },
      {
        label: "Group Receipt Settings",
        routerLink: "receipt-settings/view",
        name: "receipt-settings",
      },
    ]);
  });

  it("init tabs with group settings", () => {
    const store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      featureConfig: {
        aiPoweredReceipts: true,
      },
    });
    component.ngOnInit();

    expect(component.tabs).toEqual([
      {
        label: "Group Details",
        routerLink: "details/view",
        name: "details"
      },
      {
        label: "Group Receipt Settings",
        routerLink: "receipt-settings/view",
        name: "receipt-settings",
      },
      {
        label: "Group AI Settings",
        routerLink: "settings/view",
        name: "settings",
      },
    ]);
  });
});
