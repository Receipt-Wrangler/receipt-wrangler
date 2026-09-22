import { Component, Input } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

import { SourceFileViewerDialogComponent } from "./source-file-viewer-dialog.component";

// Stubbed rather than pulled in from SharedUiModule: the real app-image-viewer
// renders an image canvas that wants a live image load, and the point here is the
// binding, not their rendering.
@Component({ selector: "app-dialog", template: "<ng-content></ng-content>", standalone: false })
class DialogStubComponent {
  @Input() public headerText?: string;
}

@Component({ selector: "app-image-viewer", template: "", standalone: false })
class ImageViewerStubComponent {
  @Input() public imageBase64?: string;
}

@Component({ selector: "app-button", template: "", standalone: false })
class ButtonStubComponent {
  @Input() public buttonText?: string;
  @Input() public matButtonType?: string;
}

describe("SourceFileViewerDialogComponent", () => {
  let component: SourceFileViewerDialogComponent;
  let fixture: ComponentFixture<SourceFileViewerDialogComponent>;
  const dialogRef = { close: jest.fn() };
  const data = { name: "receipt.png", encodedImage: "data:image/png;base64,AAA" };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [
        SourceFileViewerDialogComponent,
        DialogStubComponent,
        ImageViewerStubComponent,
        ButtonStubComponent,
      ],
      providers: [
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(SourceFileViewerDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("renders the upload and heads the dialog with its file name", () => {
    const viewer = fixture.debugElement.query(By.directive(ImageViewerStubComponent)).componentInstance;
    const dialog = fixture.debugElement.query(By.directive(DialogStubComponent)).componentInstance;

    expect(viewer.imageBase64).toBe(data.encodedImage);
    expect(dialog.headerText).toBe(data.name);
  });

  it("closes through the dialog ref", () => {
    component.close();

    expect(dialogRef.close).toHaveBeenCalled();
  });
});
