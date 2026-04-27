package demo;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/visitor")
public class VisitorApplicationController {

    @PostMapping("/apply")
    @Operation(summary = "提交访客申请")
    public Result<VisitorApplyResponse> apply(@RequestBody VisitorApplyRequest request, @RequestHeader("X-Trace-Id") String traceId) {
        return Result.success(new VisitorApplyResponse());
    }

    @GetMapping("/{id}")
    public Result<VisitorApplyResponse> detail(@PathVariable("id") Long id, @RequestParam("verbose") Boolean verbose) {
        return Result.success(new VisitorApplyResponse());
    }
}
