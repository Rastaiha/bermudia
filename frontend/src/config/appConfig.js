// Application-level UI/behavior configuration.
//
// These are hardcoded build-time constants (mirroring the pattern used in
// `src/services/api/base_url.js`); there is no `.env`/`import.meta.env` usage
// in the frontend, so changing behavior means editing this file and rebuilding.

export const APP_CONFIG = {
    // Whether the territory map background should stick to (pan/zoom together
    // with) the islands.
    //
    //   true  — the background lives inside the panzoom'd map, so when the
    //           player drags the map the background moves with the islands.
    //           This is the intended behavior: islands stay fixed relative to
    //           the background instead of appearing to slide across a static
    //           backdrop.
    //   false — the background is rendered as a static, fixed backdrop while
    //           the islands pan/zoom on top of it (legacy behavior).
    STICK_BACKGROUND_TO_ISLANDS: true,
};
