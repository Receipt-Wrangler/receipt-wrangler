import { Observable, of } from "rxjs";

/**
 * A stand-in for the MatDialogRef that `MatDialog.open` returns, for specs that
 * stub `open()`. `setDialogInputs` writes through `componentRef.setInput`, which
 * this records in `inputs`; `componentInstance` carries whatever else the opener
 * touches on the component, such as an output it subscribes to.
 */
export function mockDialogRef(
  options: {
    afterClosed?: Observable<unknown>;
    componentInstance?: Record<string, unknown>;
  } = {}
) {
  const inputs: Record<string, any> = {};
  return {
    inputs,
    componentInstance: options.componentInstance ?? {},
    componentRef: {
      setInput: jest.fn((name: string, value: unknown) => {
        inputs[name] = value;
      }),
    },
    afterClosed: () => options.afterClosed ?? of(undefined),
    close: jest.fn(),
  };
}
