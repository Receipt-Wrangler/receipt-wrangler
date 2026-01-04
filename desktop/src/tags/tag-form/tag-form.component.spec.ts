import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef, } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { of } from "rxjs";
import { DuplicateValidator } from "src/validators/duplicate-validator";
import { ApiModule, TagService, TagView } from "../../open-api";
import { PipesModule } from "../../pipes";
import { TagFormComponent } from "./tag-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("CategoryForm", () => {
  let component: TagFormComponent;
  let fixture: ComponentFixture<TagFormComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
    declarations: [TagFormComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        MatDialogModule,
        MatSnackBarModule,
        PipesModule,
        ReactiveFormsModule],
    providers: [
        DuplicateValidator,
        {
            provide: MatDialogRef,
            useValue: {
                close: () => { },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
});
    fixture = TestBed.createComponent(TagFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form with data", () => {
    const tag: TagView = {
      id: 1,
      name: "test",
      description: "test",
      numberOfReceipts: 1,
    };
    component.tag = tag;

    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: "test",
      description: "test",
    });
  });

  it("should submit form with correct data, when editing", () => {
    const tagServiceSpy = jest.spyOn(TestBed.inject(TagService), "updateTag");
    tagServiceSpy.mockReturnValue(of({} as any));
    const nameValidateSpy = jest.spyOn(
      TestBed.inject(TagService),
      "getTagCountByName"
    ).mockReturnValue(of(0) as any);
    const tag: TagView = {
      id: 1,
      name: "test",
      description: "test",
      numberOfReceipts: 1,
    };
    component.tag = tag;

    component.ngOnInit();
    component.submit();

    expect(tagServiceSpy).toHaveBeenCalledWith(
      1,
      {
        name: "test",
        description: "test",
      },
    );
  });

  it("should submit form with correct data, when creating", () => {
    const nameValidateSpy = jest.spyOn(
      TestBed.inject(TagService),
      "getTagCountByName"
    ).mockReturnValue(of(0) as any);
    const tagServiceSpy = jest.spyOn(TestBed.inject(TagService), "createTag");
    tagServiceSpy.mockReturnValue(of({} as any));

    component.ngOnInit();
    component.form.patchValue({
      name: "test",
      description: "test",
    });
    component.submit();

    expect(tagServiceSpy).toHaveBeenCalledWith({
      name: "test",
      description: "test",
    });
  });
});
