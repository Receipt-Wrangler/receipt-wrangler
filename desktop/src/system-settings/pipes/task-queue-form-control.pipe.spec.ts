import { TaskQueueFormControlPipe } from "./task-queue-form-control.pipe";

describe("AsynqQueueFormControlPipe", () => {
  it("create an instance", () => {
    const pipe = new TaskQueueFormControlPipe();
    expect(pipe).toBeTruthy();
  });
});
