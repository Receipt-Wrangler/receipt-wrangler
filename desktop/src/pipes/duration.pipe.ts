import { Pipe, PipeTransform } from "@angular/core";
import { intervalToDuration, isToday, isValid } from "date-fns";

@Pipe({
  name: "duration",
  pure: true,
  standalone: false
})
export class DurationPipe implements PipeTransform {
  public transform(date: string | Date | null | undefined): string {
    if (!date || !isValid(new Date(date))) {
      return "";
    }

    const isDateToday = isToday(date);
    const duration = intervalToDuration({
      start: date,
      end: new Date()
    });

    const hours = duration?.hours ?? 0;
    const minutes = duration?.minutes ?? 0;
    const seconds = duration?.seconds ?? 0;
    const months = duration?.months ?? 0;
    const days = duration?.days ?? 0;

    if (days < 0 || hours < 0 || minutes < 0) {
      return "";
    }

    if (isDateToday && hours) {
      return `${hours} ${this.pluralize(hours, "hour")} ago`;
    }

    if (isDateToday && minutes) {
      return `${minutes} ${this.pluralize(minutes, "minute")} ago`;
    }

    if (isDateToday && seconds) {
      return "just now";
    }

    if (!isDateToday && months) {
      return `${months} ${this.pluralize(months, "month")} ago`;
    }

    if (!isDateToday && days) {
      return `${days} ${this.pluralize(days, "day")} ago`;
    }

    if (!isDateToday && hours) {
      return `${hours} ${this.pluralize(hours, "hour")} ago`;
    }

    return "";
  }

  private pluralize(count: number, word: string): string {
    return count === 1 ? word : word + "s";
  }
}
