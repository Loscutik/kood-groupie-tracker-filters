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

var output1 = document.getElementById("CreationYearFrom");
var output2 = document.getElementById("CreationYearTo");
output1.innerHTML = minRange.value; // Display the default slider value
minRange.oninput = function () {
  output1.innerHTML = this.value;
}
maxRange.oninput = function () {
  output2.innerHTML = this.value;
}

const selectElement = document.getElementById("frday");
const spanElement = document.getElementById("frdate");

selectElement.addEventListener("change", function () {
  const selectedOption = selectElement.options[selectElement.selectedIndex];
  spanElement.textContent = selectedOption.textContent;
});

const selectElement2 = document.getElementById("frmont");
const spanElement2 = document.getElementById("frmonth");

selectElement2.addEventListener("change", function () {
  const selectedOption2 = selectElement2.options[selectElement2.selectedIndex];
  spanElement2.textContent = selectedOption2.textContent;
});

const selectElement3 = document.getElementById("today");
const spanElement3 = document.getElementById("todate");

selectElement3.addEventListener("change", function () {
  const selectedOption3 = selectElement3.options[selectElement3.selectedIndex];
  spanElement3.textContent = selectedOption3.textContent;
});

const selectElement4 = document.getElementById("tomont");
const spanElement4 = document.getElementById("tomonth");

selectElement4.addEventListener("change", function () {
  const selectedOption4 = selectElement4.options[selectElement4.selectedIndex];
  spanElement4.textContent = selectedOption4.textContent;
});

function filter() {
  if (document.getElementById("filtr").style.display === "block") {
    document.getElementById("filtr").style.display = "none";
    window.scrollTo(0, 0);
  } else {
    document.getElementById("filtr").style.display = "block";
    window.scrollTo(0, 0);
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
slider3.oninput = function () {
  output3.innerHTML = this.value;
}
slider4.oninput = function () {
  output4.innerHTML = this.value;
}

//-------place filters from url-----

const filtersFieldsInput = [
  {
    name: "CreationDateFrom",
  },
  {
    name: "CreationDateTo",
  },
  {
    name: "FirstAlbumFromYear",
  },
  {
    name: "FirstAlbumToYear",
  },
  {
    name: "Location",
  },
  {
    name: "Name",
  },
];
const filtersFieldsSelect = [
  {
    name: "FirstAlbumFromDay",
  },
  {
    name: "FirstAlbumFromMonth",
  },
  {
    name: "FirstAlbumToDay",
  },
  {
    name: "FirstAlbumToMonth",
  },
];

const membersFields ={
  name: "Members",
  prefixId: 'm',
}

const filtersForm = document.getElementById("filters");
const searchParams = new URLSearchParams(document.location.search);

filtersFieldsInput.forEach((field) => {
  let fieldValue = searchParams.get(field.name);
  if (fieldValue) {
    const inputElm = filtersForm.querySelector("." + field.name);
    inputElm.value = fieldValue;

    const outputElm = filtersForm.querySelector(".out" + field.name)
    if (outputElm) outputElm.textContent = inputElm.value;
  }
});

filtersFieldsSelect.forEach((field) => {
  let fieldValue = searchParams.get(field.name);

  if (fieldValue && fieldValue != '0') {
    //const selectorElm = document.getElementById(field.prefixId + fieldValue);
    const selectorElm = filtersForm.querySelector(`[name="`+field.name+`"]>[value="`+fieldValue+  `"]`);
    selectorElm.setAttribute("selected", '');

    const outputElm = filtersForm.querySelector(".out" + field.name)
    //const outputElm = selectorElm.querySelector(`[value="`+fieldValue+`"]`);
    if (outputElm) outputElm.textContent = selectorElm.textContent;
  }
});

let fieldValues = searchParams.getAll(membersFields.name);
for(let number of fieldValues) {
  const memberElm = document.getElementById(membersFields.prefixId + number);
    memberElm.setAttribute("checked", '');

    
  }

