import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { GENERATED_PASSWORD_LENGTH } from "../utils/password.utils";
import {
  PASSWORD_COPIED_MESSAGE,
  PASSWORD_COPY_FAILED_MESSAGE,
  PasswordGeneratorService,
} from "./password-generator.service";
import { SnackbarService } from "./snackbar.service";

describe("PasswordGeneratorService", () => {
  let service: PasswordGeneratorService;
  let snackbarService: { success: jest.Mock; error: jest.Mock };
  let originalClipboard: any;

  beforeEach(() => {
    snackbarService = { success: jest.fn(), error: jest.fn() };

    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        { provide: SnackbarService, useValue: snackbarService },
      ],
    });

    service = TestBed.inject(PasswordGeneratorService);
    originalClipboard = (navigator as any).clipboard;
  });

  afterEach(() => {
    (navigator as any).clipboard = originalClipboard;
  });

  it("creates", () => {
    expect(service).toBeTruthy();
  });

  it("returns a generated password and copies that same value", () => {
    const copyAndNotify = jest
      .spyOn(service, "copyAndNotify")
      .mockResolvedValue(undefined);

    const password = service.generateAndCopy();

    expect(password.length).toEqual(GENERATED_PASSWORD_LENGTH);
    expect(copyAndNotify).toHaveBeenCalledWith(password);
  });

  it("writes the password to the clipboard and toasts success", async () => {
    const writeText = jest
      .fn<Promise<void>, [string]>()
      .mockResolvedValue(undefined);
    (navigator as any).clipboard = { writeText };

    await service.copyAndNotify("hunter2");

    expect(writeText).toHaveBeenCalledWith("hunter2");
    expect(snackbarService.success).toHaveBeenCalledWith(
      PASSWORD_COPIED_MESSAGE
    );
    expect(snackbarService.error).not.toHaveBeenCalled();
  });

  it("toasts an error when the clipboard write is rejected", async () => {
    (navigator as any).clipboard = {
      writeText: jest
        .fn<Promise<void>, [string]>()
        .mockRejectedValue(new Error("blocked")),
    };

    await service.copyAndNotify("hunter2");

    expect(snackbarService.error).toHaveBeenCalledWith(
      PASSWORD_COPY_FAILED_MESSAGE
    );
    expect(snackbarService.success).not.toHaveBeenCalled();
  });

  it("toasts an error without throwing when the clipboard is unavailable", async () => {
    (navigator as any).clipboard = undefined;

    await service.copyAndNotify("hunter2");

    expect(snackbarService.error).toHaveBeenCalledWith(
      PASSWORD_COPY_FAILED_MESSAGE
    );
    expect(snackbarService.success).not.toHaveBeenCalled();
  });
});
