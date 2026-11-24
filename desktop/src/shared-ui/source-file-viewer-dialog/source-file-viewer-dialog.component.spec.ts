import { Component, input } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

import { SourceFileViewerDialogComponent } from "./source-file-viewer-dialog.component";

// Stubbed rather than the real components: the real app-image-viewer
// renders an image canvas that wants a live image load, and the point here is the
// binding, not their rendering.
@Component({ selector: "app-dialog", template: "<ng-content></ng-content>" })
class DialogStubComponent {
  public readonly headerText = input<string>();
}

@Component({ selector: "app-image-viewer", template: "" })
class ImageViewerStubComponent {
  public readonly imageBase64 = input<string>();
}

@Component({ selector: "app-button", template: "" })
class ButtonStubComponent {
  public readonly buttonText = input<string>();
  public readonly matButtonType = input<string>();
}

describe("SourceFileViewerDialogComponent", () => {
  let component: SourceFileViewerDialogComponent;
  let fixture: ComponentFixture<SourceFileViewerDialogComponent>;
  const dialogRef = { close: vi.fn() };
  const data = { name: "receipt.png", encodedImage: "data:image/png;base64,AAA" };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SourceFileViewerDialogComponent],
      providers: [
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    })
      .overrideComponent(SourceFileViewerDialogComponent, {
        set: { imports: [DialogStubComponent, ImageViewerStubComponent, ButtonStubComponent] },
      })
      .compileComponents();

    fixture = TestBed.createComponent(SourceFileViewerDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("renders the upload and heads the dialog with its file name", () => {
    const viewer = fixture.debugElement.query(By.directive(ImageViewerStubComponent)).componentInstance;
    const dialog = fixture.debugElement.query(By.directive(DialogStubComponent)).componentInstance;

    expect(viewer.imageBase64()).toBe(data.encodedImage);
    expect(dialog.headerText()).toBe(data.name);
  });

  it("closes through the dialog ref", () => {
    component.close();

    expect(dialogRef.close).toHaveBeenCalled();
  });
});
