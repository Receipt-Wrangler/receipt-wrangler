import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from "@angular/material/icon";
import { MatListModule } from "@angular/material/list";
import { InputModule } from "../../input/index";
import { PipesModule } from "../../pipes/index";
import { EditableListComponent } from "../../shared-ui/editable-list/editable-list.component";
import { IconAutocompleteComponent } from "../../shared-ui/icon-autocomplete/icon-autocomplete.component";
import { UserShortcutComponent } from "./user-shortcut.component";

describe("UserShortcutComponent", () => {
  let component: UserShortcutComponent;
  let fixture: ComponentFixture<UserShortcutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [UserShortcutComponent, EditableListComponent, IconAutocompleteComponent],
      imports: [MatIconModule, InputModule, PipesModule, MatListModule, ReactiveFormsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(UserShortcutComponent);
    component = fixture.componentInstance;
    component.parentForm = new FormGroup({
      userShortcuts: new FormArray([])
    });
    component.formConfig = {} as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should emit formCommand event when removeShortcut is called", () => {
    jest.spyOn(component.formCommand, "emit");
    component.removeShortcut(1);
    expect(component.formCommand.emit).toHaveBeenCalledWith({
      path: "userShortcuts",
      command: "removeAt",
      payload: 1
    });
  });

  it("should emit shortcutDoneClicked event when emitShortcutDoneClicked is called", () => {
    jest.spyOn(component.shortcutDoneClicked, "emit");
    component.emitShortcutDoneClicked();
    expect(component.shortcutDoneClicked.emit).toHaveBeenCalled();
  });

  it("should emit shortcutCancelClicked event when emitShortcutCancelClicked is called", () => {
    jest.spyOn(component.shortcutCancelClicked, "emit");
    component.emitShortcutCancelClicked();
    expect(component.shortcutCancelClicked.emit).toHaveBeenCalled();
  });

  it("should return userShortcuts form array", () => {
    expect(component.userShortcuts).toBeInstanceOf(FormArray);
  });
});
