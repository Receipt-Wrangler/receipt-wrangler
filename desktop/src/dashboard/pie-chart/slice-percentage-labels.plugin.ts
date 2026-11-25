import { ArcElement, Plugin } from "chart.js";
import { toFont } from "chart.js/helpers";

/** A slice at or below this share of the pie gets no label of its own: it would not fit. */
const MIN_LABELLED_PERCENTAGE = 5;

/** The size of the whole pie. It counts every slice, including any hidden through the legend. */
export function sliceTotal(data: number[]): number {
  return data.reduce((sum, slice) => sum + Math.abs(slice), 0);
}

/** A slice's share of [total] with one decimal ("20.0"), or "0" when the pie is empty. */
export function slicePercentage(value: number, total: number): string {
  return total > 0 ? ((Math.abs(value) / total) * 100).toFixed(1) : "0";
}

/** A slice's share of the pie made of [data]. */
export function getSlicePercentage(value: number, data: number[]): string {
  return slicePercentage(value, sliceTotal(data));
}

/** Writes each visible slice's percentage, in bold white, at the middle of the slice. */
export const slicePercentageLabelsPlugin: Plugin<"pie"> = {
  id: "slicePercentageLabels",
  afterDatasetsDraw(chart) {
    const { ctx } = chart;
    ctx.save();
    ctx.fillStyle = "#fff";
    ctx.font = toFont({ weight: "bold", size: 12 }).string;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";

    chart.data.datasets.forEach((dataset, datasetIndex) => {
      if (!chart.isDatasetVisible(datasetIndex)) {
        return;
      }
      // Once per dataset: this hook runs on every animation frame.
      const total = sliceTotal(dataset.data);
      chart.getDatasetMeta(datasetIndex).data.forEach((slice, index) => {
        if (!chart.getDataVisibility(index)) {
          return;
        }
        const percentage = slicePercentage(dataset.data[index], total);
        if (parseFloat(percentage) > MIN_LABELLED_PERCENTAGE) {
          const { x, y } = (slice as ArcElement).getCenterPoint(false);
          ctx.fillText(`${percentage}%`, x, y);
        }
      });
    });

    ctx.restore();
  },
};
