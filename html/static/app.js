function showTime(){
  let now = new Date();

  let day = now.getDate();
  let month = now.getMonth();

  let h = now.getHours(); // 0 - 23
  let m = now.getMinutes(); // 0 - 59
  let s = now.getSeconds(); // 0 - 59

  h = (h < 10) ? "0" + h : h;
  m = (m < 10) ? "0" + m : m;
  s = (s < 10) ? "0" + s : s;

  let time = `${day}/${month} ${h}:${m}:${s}`;

  document.getElementById("clock").innerHTML = time;

  setTimeout(showTime, 1000);
}

showTime();
