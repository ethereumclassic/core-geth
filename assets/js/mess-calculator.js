// MESS (ECBP-1100) confirmation calculator. Follows ecbp1100PolynomialV in core/blockchain_af.go.
(function () {
  var DENOMINATOR = 128n, XCAP = 25132n, HEIGHT = 3840n;

  function numerator(x) {
    if (x > XCAP) x = XCAP;
    return DENOMINATOR + ((3n * x * x - (2n * x * x * x) / XCAP) * HEIGHT) / (XCAP * XCAP);
  }

  function setup() {
    var form = document.getElementById("mess-calculator");
    if (!form) return;
    var hours = form.elements.namedItem("hours");
    var out = form.elements.namedItem("ratio");
    function update() {
      var seconds = Math.floor((Number(hours.value) || 0) * 3600);
      if (seconds < 0) seconds = 0;
      out.value = (Number(numerator(BigInt(seconds))) / 128).toFixed(2);
    }
    form.addEventListener("input", update);
    update();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", setup);
  } else {
    setup();
  }
})();
