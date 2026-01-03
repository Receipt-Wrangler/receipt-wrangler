import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatListModule } from "@angular/material/list";
import { ButtonModule } from "../../button/index";
import { EditableListComponent } from "./editable-list.component";

describe("EditableListComponent", () => {
  let component: EditableListComponent;
  let fixture: ComponentFixture<EditableListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [EditableListComponent],
      imports: [MatListModule, ButtonModule],
    }).compileComponents();

    fixture = TestBed.createComponent(EditableListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should emit editButtonClicked event when handleEditButtonClicked is called", () => {
    jest.spyOn(component.editButtonClicked, "emit");
    component.handleEditButtonClicked(1);
    expect(component.editButtonClicked.emit).toHaveBeenCalledWith(1);
  });

  it("should emit deleteButtonClicked event when handleDeleteButtonClicked is called", () => {
    jest.spyOn(component.deleteButtonClicked, "emit");
    component.handleDeleteButtonClicked(1);
    expect(component.deleteButtonClicked.emit).toHaveBeenCalledWith(1);
  });

  it("should set rowOpen to the correct index when handleEditButtonClicked is called", () => {
    component.handleEditButtonClicked(1);
    expect(component.getCurrentRowOpen()).toBe(1);
  });

  it("should set rowOpen to undefined when handleDeleteButtonClicked is called", () => {
    component.handleDeleteButtonClicked(1);
    expect(component.getCurrentRowOpen()).toBeUndefined();
  });

  it("should open the last row when openLastRow is called", () => {
    component.listData = [{}, {}, {}];
    component.openLastRow();
    expect(component.getCurrentRowOpen()).toBe(2);
  });

  it("should close the row when closeRow is called", () => {
    component.handleEditButtonClicked(1);
    component.closeRow();
    expect(component.getCurrentRowOpen()).toBeUndefined();
  });
});
