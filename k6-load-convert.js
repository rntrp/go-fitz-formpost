import { check } from "k6";
import http from "k6/http";

const url = "http://localhost:8080/convert?width=256&height=256&format=jpeg";
const doc = open("internal/rest/test.pdf", "b");

export default function () {
    const res = http.post(url, {
        doc: http.file(doc, "test"),
    });
    check(res, {
        http200: http.expectedStatuses(200),
    });
}
