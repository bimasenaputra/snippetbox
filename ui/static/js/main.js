function init() {
	var navbar = document.getElementsByTagName("nav")[0];
	var navPos = navbar.offsetTop;
	let navLinks = document.querySelectorAll("nav a");
	
	for (let i = 0; i < navLinks.length; i++) {
		var link = navLinks[i]
		if (link.getAttribute('href') == window.location.pathname) {
			link.classList.add("live");
			break;
		}
	}

	window.addEventListener('scroll', () => {
		if (window.scrollY >= navPos) {
			navbar.classList.add("sticky");
		} else {
			navbar.classList.remove("sticky");
		}
	})
}

window.onload = init