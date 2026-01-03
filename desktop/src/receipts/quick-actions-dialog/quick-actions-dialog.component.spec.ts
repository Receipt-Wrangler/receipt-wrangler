import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormControl, FormGroup, ReactiveFormsModule, } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { User } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services/snackbar.service";
import { QuickActionsDialogComponent } from "./quick-actions-dialog.component";

describe("QuickActionsDialogComponent", () => {
  let component: QuickActionsDialogComponent;
  let fixture: ComponentFixture<QuickActionsDialogComponent>;

  const mockDialogRef = {
    close: jest.fn(),
  };

  const mockSnackBar = {
    open: jest.fn(),
  };

  const mockSnackbarService = {
    error: jest.fn(),
    info: jest.fn(),
    success: jest.fn(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [QuickActionsDialogComponent],
      imports: [NgxsModule.forRoot([]), PipesModule, ReactiveFormsModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: MatDialogRef, useValue: mockDialogRef },
        { provide: MatSnackBar, useValue: mockSnackBar },
        { provide: SnackbarService, useValue: mockSnackbarService },
        { provide: ActivatedRoute, useValue: {} },
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QuickActionsDialogComponent);
    component = fixture.componentInstance;
    component.originalReceipt = { id: 1 } as any;
    component.amountToSplit = 100; // Set the amount to split
    fixture.detectChanges();

    // Reset mocks
    mockDialogRef.close.mockClear();
    mockSnackbarService.error.mockClear();
    mockSnackbarService.info.mockClear();
    mockSnackbarService.success.mockClear();
  });

  // Helper function to create test users
  const createTestUsers = (): User[] => [
    { id: 1, displayName: "John Doe" } as User,
    { id: 2, displayName: "Jane Smith" } as User,
    { id: 3, displayName: "Bob Johnson" } as User,
  ];

  // Helper function to setup users in the form
  const setupUsersInForm = (users: User[]): void => {
    const formArray = component.localForm.get("usersToSplit") as FormArray;
    users.forEach((user) => {
      formArray.push(new FormControl(user));
    });
    // Trigger the subscription manually
    formArray.updateValueAndValidity();
    formArray.setValue(users);
  };

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize the form on ngOnInit", () => {
    component.ngOnInit();
    expect(component.localForm).toBeDefined();
  });

  it("should add even split items", () => {
    component.amountToSplit = 100;
    component.ngOnInit();

    const users = [
      { id: 1, displayName: "User 1" },
      { id: 2, displayName: "User 2" },
    ];

    const formArray = component.localForm.get("usersToSplit") as FormArray;
    users.forEach((user) => {
      formArray.push(new FormControl(user));
    });
    formArray.setValue(users);
    component.localForm.patchValue({
      quickAction: component.radioValues[0].value,
    });

    jest.spyOn(component.itemsToAdd, "emit");

    component.addSplits();

    expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
      items: expect.arrayContaining([
        expect.objectContaining({ amount: 50, chargedToUserId: 1 }),
        expect.objectContaining({ amount: 50, chargedToUserId: 2 }),
      ]),
      itemIndex: undefined
    });
  });

  it("should split evenly with optional parts", () => {
    component.amountToSplit = 100;
    component.ngOnInit();

    const formArray = component.localForm.get("usersToSplit") as FormArray;
    const users = [
      { id: 1, displayName: "User 1" },
      { id: 2, displayName: "User 2" },
    ];
    users.forEach((user) => {
      formArray.push(new FormControl(user));
    });
    formArray.setValue(users);
    component.localForm.patchValue({
      quickAction: component.radioValues[1].value,
    });
    component.localForm.addControl("1", new FormControl("10"));
    component.localForm.addControl("2", new FormControl("20"));

    jest.spyOn(component.itemsToAdd, "emit");

    component.addSplits();

    expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
      items: expect.arrayContaining([
        expect.objectContaining({ amount: 1, chargedToUserId: 1, name: "User 1's Portion" }),
        expect.objectContaining({ amount: 1, chargedToUserId: 2, name: "User 2's Portion" }),
        expect.objectContaining({ amount: 49, chargedToUserId: 1, name: "User 1's Even Portion" }),
        expect.objectContaining({ amount: 49, chargedToUserId: 2, name: "User 2's Even Portion" }),
      ]),
      itemIndex: undefined
    });
  });

  describe("Split by Percentage", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should include Split by Percentage in radio options", () => {
      expect(component.radioValues).toContainEqual(
        expect.objectContaining({
          displayValue: "Split by Percentage",
          value: "SplitByPercentage",
        })
      );
    });

    it("should create percentage controls when users are added", () => {
      const users = createTestUsers();
      setupUsersInForm(users);

      users.forEach((user) => {
        expect(component.localForm.get(`${user.id}_percentage`)).toBeDefined();
        expect(component.localForm.get(`${user.id}_customPercentage`)).toBeDefined();
      });
    });

    it("should initialize percentage controls as disabled", () => {
      const users = createTestUsers();
      setupUsersInForm(users);

      users.forEach((user) => {
        const percentageControl = component.localForm.get(`${user.id}_percentage`);
        expect(percentageControl?.disabled).toBe(true);
      });
    });

    it("should enable percentage control when custom checkbox is checked", () => {
      const users = createTestUsers();
      setupUsersInForm(users);

      const user = users[0];
      const customControl = component.localForm.get(`${user.id}_customPercentage`);
      const percentageControl = component.localForm.get(`${user.id}_percentage`);

      customControl?.setValue(true);

      expect(percentageControl?.enabled).toBe(true);
    });

    it("should disable percentage control and reset value when custom checkbox is unchecked", () => {
      const users = createTestUsers();
      setupUsersInForm(users);

      const user = users[0];
      const customControl = component.localForm.get(`${user.id}_customPercentage`);
      const percentageControl = component.localForm.get(`${user.id}_percentage`);

      // Enable and set value
      customControl?.setValue(true);
      percentageControl?.setValue(50);

      // Disable
      customControl?.setValue(false);

      expect(percentageControl?.disabled).toBe(true);
      expect(percentageControl?.value).toBe(0);
    });

    it("should set predefined percentage when button is clicked", () => {
      const users = createTestUsers();
      setupUsersInForm(users);

      const user = users[0];
      component.setPercentage(user.id.toString(), 75);

      const percentageControl = component.localForm.get(`${user.id}_percentage`);
      const customControl = component.localForm.get(`${user.id}_customPercentage`);

      expect(percentageControl?.value).toBe(75);
      expect(customControl?.value).toBe(false);
      expect(percentageControl?.disabled).toBe(true);
    });

    describe("Validation", () => {
      it("should validate total percentage not exceeding 100%", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Set percentages that exceed 100%
        component.setPercentage(users[0].id.toString(), 60);
        component.setPercentage(users[1].id.toString(), 50);

        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          "Total percentage cannot exceed 100!"
        );
      });

      it("should validate total percentage is greater than 0%", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // All percentages remain at 0
        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          "Total percentage must be greater than 0!"
        );
      });

      it("should validate individual percentage is between 0 and 100", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Set invalid percentage
        const percentageControl = component.localForm.get(`${users[0].id}_percentage`);
        percentageControl?.enable();
        percentageControl?.setValue(150);

        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          `Percentage for ${users[0].displayName} must be between 0 and 100!`
        );
      });

      it("should validate custom percentage field is not empty when custom is enabled", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Enable custom but leave empty
        const customControl = component.localForm.get(`${users[0].id}_customPercentage`);
        const percentageControl = component.localForm.get(`${users[0].id}_percentage`);

        customControl?.setValue(true);
        percentageControl?.setValue("");

        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          `Please enter a percentage for ${users[0].displayName}!`
        );
      });

      it("should pass validation with valid percentages", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Set valid percentages
        component.setPercentage(users[0].id.toString(), 40);
        component.setPercentage(users[1].id.toString(), 35);
        component.setPercentage(users[2].id.toString(), 25);

        component.addSplits();

        expect(mockDialogRef.close).toHaveBeenCalledWith(true);
      });
    });

    describe("Split Calculation", () => {
      it("should create items with correct amounts based on percentages", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Set percentages: 40%, 35%, 25%
        component.setPercentage(users[0].id.toString(), 40);
        component.setPercentage(users[1].id.toString(), 35);
        component.setPercentage(users[2].id.toString(), 25);

        jest.spyOn(component.itemsToAdd, "emit");

        component.addSplits();

        expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
          items: expect.arrayContaining([
            expect.objectContaining({ amount: 40 }), // 100 * 40% = 40
            expect.objectContaining({ amount: 35 }), // 100 * 35% = 35
            expect.objectContaining({ amount: 25 }), // 100 * 25% = 25
          ]),
          itemIndex: undefined
        });
      });

      it("should create items with correct names including percentage", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        component.setPercentage(users[0].id.toString(), 50);

        jest.spyOn(component.itemsToAdd, "emit");

        component.addSplits();

        expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
          items: expect.arrayContaining([
            expect.objectContaining({ name: "John Doe's 50% Portion" }),
          ]),
          itemIndex: undefined
        });
      });

      it("should only create items for users with percentage > 0", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Only set percentage for first user
        component.setPercentage(users[0].id.toString(), 100);
        // Other users remain at 0%

        jest.spyOn(component.itemsToAdd, "emit");

        component.addSplits();

        expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
          items: [
            expect.objectContaining({ chargedToUserId: users[0].id }),
          ],
          itemIndex: undefined
        });
      });

      it("should handle decimal percentages correctly", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        // Set decimal percentage
        const percentageControl = component.localForm.get(`${users[0].id}_percentage`);
        percentageControl?.enable();
        percentageControl?.setValue(33.33);

        jest.spyOn(component.itemsToAdd, "emit");

        component.addSplits();

        expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
          items: [
            expect.objectContaining({ amount: 33.33 }),
          ],
          itemIndex: undefined
        });
      });
    });

    describe("Form State Management", () => {
      it("should clear percentage errors when switching away from percentage mode", () => {
        component.ngOnInit();
        const users = createTestUsers();
        setupUsersInForm(users);

        // Set percentage mode and create some errors
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const percentageControl = component.localForm.get(`${users[0].id}_percentage`);
        percentageControl?.markAsTouched();
        percentageControl?.setErrors({ invalid: true });

        // Switch to even split
        component.localForm.patchValue({
          quickAction: component.radioValues[0].value,
        });

        expect(percentageControl?.errors).toBe(null);
        expect(percentageControl?.touched).toBe(false);
      });

      it("should update validators when custom percentage is enabled", () => {
        const users = createTestUsers();
        setupUsersInForm(users);

        const user = users[0];
        const customControl = component.localForm.get(`${user.id}_customPercentage`);
        const percentageControl = component.localForm.get(`${user.id}_percentage`);

        // Enable custom percentage
        customControl?.setValue(true);

        // Check if required validator is added
        percentageControl?.setValue("");
        percentageControl?.updateValueAndValidity();

        expect(percentageControl?.hasError("required")).toBe(true);
      });

      it("should remove controls when users are removed", () => {
        let users = createTestUsers();
        setupUsersInForm(users);

        // Verify controls exist for all users
        users.forEach(user => {
          expect(component.localForm.get(`${user.id}_percentage`)).toBeDefined();
          expect(component.localForm.get(`${user.id}_customPercentage`)).toBeDefined();
        });

        // Remove a user by properly updating the FormArray
        users = users.slice(1); // Remove first user
        const formArray = component.localForm.get("usersToSplit") as FormArray;
        
        // Clear the FormArray and rebuild it with remaining users
        formArray.clear();
        users.forEach((user) => {
          formArray.push(new FormControl(user));
        });
        
        // Trigger the valueChanges subscription manually
        formArray.updateValueAndValidity();

        // Verify control for removed user is gone
        expect(component.localForm.get(`1_percentage`)).toBe(null);
        expect(component.localForm.get(`1_customPercentage`)).toBe(null);
        
        // Verify controls for remaining users still exist
        users.forEach(user => {
          expect(component.localForm.get(`${user.id}_percentage`)).toBeDefined();
          expect(component.localForm.get(`${user.id}_customPercentage`)).toBeDefined();
        });
      });
    });

    describe("Edge Cases", () => {
      it("should handle receipt amount of 0", () => {
        component.amountToSplit = 0;
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          "Amount to split does not exist or is invalid!"
        );
      });

      it("should handle invalid receipt amount", () => {
        component.amountToSplit = NaN;
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        component.addSplits();

        expect(mockSnackbarService.error).toHaveBeenCalledWith(
          "Amount to split does not exist or is invalid!"
        );
      });

      it("should handle percentage values with more than 2 decimal places", () => {
        component.localForm.patchValue({
          quickAction: component.radioValues[2].value,
        });

        const users = createTestUsers();
        setupUsersInForm(users);

        const percentageControl = component.localForm.get(`${users[0].id}_percentage`);
        percentageControl?.enable();
        percentageControl?.setValue(33.333);

        jest.spyOn(component.itemsToAdd, "emit");

        component.addSplits();

        // Should round to 2 decimal places
        expect(component.itemsToAdd.emit).toHaveBeenCalledWith({
          items: [
            expect.objectContaining({ amount: 33.33 }),
          ],
          itemIndex: undefined
        });
      });
    });
  });
});
