var slider = document.getElementById("myRange");
var output = document.getElementById("CreationYear");
output.innerHTML = slider.value; // Display the default slider value

// Update the current slider value (each time you drag the slider handle)
slider.oninput = function() {
  output.innerHTML = this.value;
}

var slider1 = document.getElementById("myRange2");
var output1 = document.getElementById("CreationYear22");
output1.innerHTML = slider1.value; // Display the default slider value

// Update the current slider value (each time you drag the slider handle)
slider1.oninput = function() {
  output1.innerHTML = this.value;
}

//-----------------
function filter() {
  if (document.getElementById("filter").style.display == "none") {
      document.getElementById("filter").style.display = "block";
  } else {
      document.getElementById("filter").style.display = "none";
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