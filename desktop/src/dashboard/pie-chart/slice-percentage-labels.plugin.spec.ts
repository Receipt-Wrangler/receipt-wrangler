import type { Chart } from "chart.js";
import { getSlicePercentage, slicePercentageLabelsPlugin } from "./slice-percentage-labels.plugin";

/**
 * One pie dataset of [data]. Slice i's center point is (i * 10, i * 20), so a label's position
 * tells which slice it was drawn for.
 */
function fakePie(data: number[], options: { hidden?: number[]; datasetVisible?: boolean } = {}) {
  const ctx = {
    save: vi.fn(),
    restore: vi.fn(),
    fillText: vi.fn(),
    fillStyle: "",
    font: "",
    textAlign: "",
    textBaseline: "",
  };
  const chart = {
    ctx,
    data: { datasets: [{ data }] },
    isDatasetVisible: () => options.datasetVisible ?? true,
    getDataVisibility: (index: number) => !(options.hidden ?? []).includes(index),
    getDatasetMeta: () => ({
      data: data.map((_, index) => ({ getCenterPoint: () => ({ x: index * 10, y: index * 20 }) })),
    }),
  };
  return { chart: chart as unknown as Chart<"pie">, ctx };
}

function draw(chart: Chart<"pie">): void {
  slicePercentageLabelsPlugin.afterDatasetsDraw!(chart, {}, {}, false);
}

describe("slicePercentageLabelsPlugin", () => {
  it("writes each slice's share at the middle of the slice", () => {
    const { chart, ctx } = fakePie([100, 200, 200]);

    draw(chart);

    expect(ctx.fillText.mock.calls).toEqual([
      ["20.0%", 0, 0],
      ["40.0%", 10, 20],
      ["40.0%", 20, 40],
    ]);
  });

  it("draws in bold 12px white, centered, and restores the canvas state", () => {
    const { chart, ctx } = fakePie([100, 200]);

    draw(chart);

    expect(ctx.fillStyle).toBe("#fff");
    expect(ctx.font).toMatch(/^normal bold 12px /);
    expect(ctx.textAlign).toBe("center");
    expect(ctx.textBaseline).toBe("middle");
    expect(ctx.save).toHaveBeenCalledTimes(1);
    expect(ctx.restore).toHaveBeenCalledTimes(1);
  });

  // The threshold applies to the rounded figure, so 5.04% reads "5.0" and stays unlabelled.
  it("leaves a slice of 5% or less unlabelled", () => {
    const { chart, ctx } = fakePie([504, 506, 8990]);

    draw(chart);

    expect(ctx.fillText.mock.calls.map(([text]) => text)).toEqual(["5.1%", "89.9%"]);
  });

  // The total still counts the hidden slice, so the visible slices keep their figures.
  it("skips a slice hidden through the legend", () => {
    const { chart, ctx } = fakePie([100, 200, 200], { hidden: [1] });

    draw(chart);

    expect(ctx.fillText.mock.calls).toEqual([
      ["20.0%", 0, 0],
      ["40.0%", 20, 40],
    ]);
  });

  it("skips a hidden dataset", () => {
    const { chart, ctx } = fakePie([100, 200], { datasetVisible: false });

    draw(chart);

    expect(ctx.fillText).not.toHaveBeenCalled();
  });

  it("draws nothing for an empty pie", () => {
    const { chart, ctx } = fakePie([0, 0]);

    draw(chart);

    expect(ctx.fillText).not.toHaveBeenCalled();
  });
});

describe("getSlicePercentage", () => {
  it("gives the slice's share of the total to one decimal", () => {
    expect(getSlicePercentage(100, [100, 200, 200])).toBe("20.0");
    expect(getSlicePercentage(100, [100, 200])).toBe("33.3");
  });

  it("measures negative values by their size", () => {
    expect(getSlicePercentage(-100, [-100, 300])).toBe("25.0");
  });

  it("gives 0 when the total is zero", () => {
    expect(getSlicePercentage(0, [0, 0, 0])).toBe("0");
  });
});
