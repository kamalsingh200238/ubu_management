import htmx from "htmx.org";
import _hyperscript from "hyperscript.org";
import "./shoelace.ext.js";
import "@shoelace-style/shoelace";
import "@shoelace-style/shoelace/dist/themes/light.css";
import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
setBasePath("");

window.htmx = htmx;
_hyperscript.browserInit();

document.addEventListener("alert", (event) => {
  console.log({ event });

  const container = document.getElementById("alert-toast-wrapper");
  let icon;

  // Select icon based on variant
  switch (event.detail.variant) {
    case "primary":
      icon = "info-circle";
      break;
    case "success":
      icon = "check2-circle";
      break;
    case "neutral":
      icon = "gear";
      break;
    case "warning":
      icon = "exclamation-triangle";
      break;
    case "danger":
      icon = "exclamation-octagon";
      break;
    default:
      icon = "info-circle";
  }

  const alert = Object.assign(document.createElement("sl-alert"), {
    variant: event.detail.variant,
    closable: event.detail.closable,
    duration: event.detail.duration,
    innerHTML: `
      <sl-icon name="${icon}" slot="icon"></sl-icon>
      ${event.detail.message}
    `,
  });

  // Append the alert to the alert-toast-wrapper container
  container.append(alert);
  alert.toast();
});
