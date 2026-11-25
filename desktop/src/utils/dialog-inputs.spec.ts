import { ApplicationRef, Component, input } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { firstValueFrom } from "rxjs";
import { setDialogInputs } from "./dialog-inputs";

@Component({
  template: `{{ headerText() }}: {{ count() }}`,
})
class DialogContentComponent {
  public readonly headerText = input<string>("");

  public readonly count = input.required<number>();
}

describe("setDialogInputs", () => {
  afterEach(() => TestBed.inject(MatDialog).closeAll());

  it("sets the component's inputs before its first render", async () => {
    const dialogRef = TestBed.inject(MatDialog).open(DialogContentComponent);

    setDialogInputs(dialogRef, { headerText: "Delete tag", count: 3 });
    await firstValueFrom(dialogRef.afterOpened());
    // With animations off (src/test-providers.ts) the dialog reports itself open before the
    // zoneless scheduler has rendered it, so let the app settle before reading the DOM.
    await TestBed.inject(ApplicationRef).whenStable();

    expect(dialogRef.componentInstance.headerText()).toEqual("Delete tag");
    expect(dialogRef.componentInstance.count()).toEqual(3);
    expect(document.body.textContent).toContain("Delete tag: 3");
  });

  it("throws for a dialog opened from a template", () => {
    const dialogRef = { componentRef: null } as MatDialogRef<DialogContentComponent>;

    expect(() => setDialogInputs(dialogRef, { headerText: "x" })).toThrow(
      /opened with a component/
    );
  });
});
