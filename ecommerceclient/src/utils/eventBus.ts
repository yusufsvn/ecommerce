export const requestStartEvent = new Event('api-request-start');
export const requestEndEvent = new Event('api-request-end');

export function triggerLoader(show: boolean) {
    if (show) {
        window.dispatchEvent(requestStartEvent);
    } else {
        window.dispatchEvent(requestEndEvent);
    }
}
