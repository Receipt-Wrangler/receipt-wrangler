import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { provideTestStore } from "../../store/test-store";

import { UserAutocompleteComponent } from "./user-autocomplete.component";

describe("UserAutocompleteComponent", () => {
  let component: UserAutocompleteComponent;
  let fixture: ComponentFixture<UserAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserAutocompleteComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideTestStore(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(UserAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('inputFormControl', new FormControl());
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should not clear inputFormControl value on initial load when groupId is not set", async () => {
    const control = new FormControl("existing-user-id");
    fixture.componentRef.setInput('inputFormControl', control);
    await fixture.whenStable();

    expect(control.value).toBe("existing-user-id");
  });
});
