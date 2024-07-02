import htmx from "htmx.org";
import _hyperscript from "hyperscript.org";

import "./shoelace.ext.js";
import "@shoelace-style/shoelace";
import "@shoelace-style/shoelace/dist/themes/light.css";
import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
import { createShoelaceAlert } from "./utils/create-shoelace-alert.js";
import { setAnimation } from "@shoelace-style/shoelace/dist/utilities/animation-registry.js";
setBasePath("");


window.htmx = htmx;
_hyperscript.browserInit();

setAnimation('dialog.show', {
  options: {
    duration: 50000,
  },
});

document.addEventListener("alert", (event) => {
  createShoelaceAlert({
    variant: event.detail.variant,
    closable: event.detail.closable,
    duration: event.detail.duration,
    heading: event.detail.heading,
    messages: event.detail.messages,
  });
});
