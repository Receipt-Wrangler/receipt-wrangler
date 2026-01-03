import { animate, state, style, transition, trigger, } from "@angular/animations";

export const fadeInOut = [
  trigger("fadeInOut", [
    state(
      "void",
      style({
        opacity: 0,
        visibility: "hidden",
      })
    ),
    transition(":enter", [
      animate(
        "0.2s",
        style({
          opacity: 1,
          visibility: "visible",
        })
      ),
    ]),
    transition(":leave", [
      animate(
        "0.2s",
        style({
          opacity: 0,
          visibility: "hidden",
        })
      ),
    ]),
  ]),
];

export const fadeIn = [
  trigger("fadeIn", [
    state(
      "void",
      style({
        opacity: 0,
        visibility: "hidden",
      })
    ),
    transition(":enter", [
      animate(
        "0.2s",
        style({
          opacity: 1,
          visibility: "visible",
        })
      ),
    ]),
  ]),
];
