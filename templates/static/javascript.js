const minRange = document.getElementById('minRange');
const maxRange = document.getElementById('maxRange');

minRange.addEventListener('input', updateSlider);
maxRange.addEventListener('input', updateSlider);

function updateSlider() {
  const min = parseInt(minRange.value);
  const max = parseInt(maxRange.value);

  if (min > max) {
    minRange.value = max;
  } else if (max < min) {
    maxRange.value = min;
  }
}

var slider1 = document.getElementById("minRange");
var slider2 = document.getElementById("maxRange");
var output1 = document.getElementById("CreationYearFrom");
var output2 = document.getElementById("CreationYearTo");
output1.innerHTML = slider1.value; // Display the default slider value
// Update the current slider value (each time you drag the slider handle)
slider1.oninput = function() {
  output1.innerHTML = this.value;
}
slider2.oninput = function() {
  output2.innerHTML = this.value;
}
//-----------------
function filter() {
  if (document.getElementById("filtr").style.display == "none") {
      window.scrollTo(0, 0);
      document.getElementById("filtr").style.display = "block";
  } else {
      window.scrollTo(0, 0);
      document.getElementById("filtr").style.display = "none";
  }
}
function rangeValue() {
  const value = document.querySelectorAll(".outYear")
  const input = document.querySelectorAll(".inYear")
  for (let index = 0; index < value.length; index++) {
      value[index].textContent = input[index].value
      input[index].addEventListener("input", (event) => {
          value[index].textContent = event.target.value
      })
  }
}

var slider3 = document.getElementById("minnYear");
var output3 = document.getElementById("fryear");
var slider4 = document.getElementById("maxxYear");
var output4 = document.getElementById("toyear");
//output3.innerHTML = slider3.value; // Display the default slider value
//output4.innerHTML = slider4.value;
// Update the current slider value (each time you drag the slider handle)
slider3.oninput = function() {
  output3.innerHTML = this.value;
}
slider4.oninput = function() {
  output4.innerHTML = this.value;
}