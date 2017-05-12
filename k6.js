import http from "k6/http";
import { check } from "k6";

const duration = "5s",
    vus = 10,
    httpStatusOK = 200,
    maxDuration = 20

export let options = {
    vus: vus,
    duration: duration
}

export default function () {
    let metrics = http.get("http://localhost:4222/api/metrics/redis")
    check(metrics, {
        "status was 200": (r) => r.status === httpStatusOK,
        "transaction time OK": (r) => r.timings.duration < maxDuration
    })

    let logs = http.get("http://localhost:4222/api/logs/redis")
    check(logs, {
        "status was 200": (r) => r.status === httpStatusOK,
        "transaction time OK": (r) => r.timings.duration < maxDuration
    })

    let stopped = http.get("http://localhost:4222/api/stopped")
    check(stopped, {
        "status was 200": (r) => r.status === httpStatusOK,
        "transaction time OK": (r) => r.timings.duration < maxDuration
    })

    let launched = http.get("http://localhost:4222/api/launched")
    check(launched, {
        "status was 200": (r) => r.status === httpStatusOK,
        "transaction time OK": (r) => r.timings.duration < maxDuration
    })
}
