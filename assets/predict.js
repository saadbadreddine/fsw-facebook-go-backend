window.onload = () => {
	class Person {
		constructor(name, age, gender, gender_probability, country, country_probability) {
			this.name = name;
			this.age = age;
			this.gender = gender;
			this.gender_probability = gender_probability;
			this.country = country;
			this.country_probability = country_probability;
		}
	}

	const fetchDog = async () => {
		const response = await fetch('https://dog.ceo/api/breeds/image/random');
		const doggo = await response.json();
		return doggo;
	};
	fetchDog().then((doggo) => {
		document.getElementById('doggo').src = doggo.message;
		document.getElementById('doggo').style.display = 'block';
	});

	/*
	fetch('https://dog.ceo/api/breeds/image/random').then((response) => response.json()).then((data) => {
		document.getElementById('doggo').src = data.message;
		document.getElementById('doggo').style.display = 'block';
	});*/

	let person = new Person();

	document.getElementById('name-form').addEventListener('submit', (e) => {
		person.name = document.getElementById('name-input').value;
		e.preventDefault();
		document.getElementById('name-form').reset();

		if (person.name === '') {
			document.getElementsByClassName('predict-box')[0].style.display = 'none';
			return;
		}

		fetchAPIs(person.name)
			.then(([ age, gender, country ]) => {
				person.name = person.name.charAt(0).toUpperCase() + person.name.slice(1).toLowerCase();
				person.age = age.age;
				person.gender = gender.gender;
				person.gender = person.gender.charAt(0).toUpperCase() + person.gender.slice(1);
				person.country = [];
				country_codes = '';
				country_names = '';

				document.getElementById('name').textContent = ` ${person.name}`;
				document.getElementById('gender').textContent = ` ${person.gender}`;
				document.getElementById('age').textContent = ` ${person.age}`;
				country.country.forEach((element) => person.country.push(element.country_id));

				person.country.forEach((element) => (country_codes += `${element},`));
				//country_codes = country_codes.replace(/,\s*$/, '');
				fetchCountries(country_codes)
					.then((data) => {
						data.forEach((element) => {
							country_names += `${element.demonyms.eng.f}, `;
						});
						country_names = country_names.replace(/,\s*$/, '');
						document.getElementById('nationality').textContent = ` ${country_names}`;
						document.getElementsByClassName('predict-box')[0].style.display = 'flex';
					})
					.catch((error) => {
						console.log(error);
					});
			})
			.catch((error) => {
				// if there's an error, log it
				alert(error);
				document.getElementsByClassName('predict-box')[0].style.display = 'none';
				console.log('Error');
				console.log(error);
			});
	});

	//-------------------------------------async/await & Promise.all-----------------------------------------//

	const fetchAPIs = async (name) => {
		const [ response_age, response_gender, response_country ] = await Promise.all([
			fetch(`https://api.agify.io/?name=${name}`),
			fetch(`https://api.genderize.io?name=${name}`),
			fetch(`https://api.nationalize.io/?name=${name}`)
		]);

		const age = await response_age.json();
		const gender = await response_gender.json();
		const country = await response_country.json();

		return [ age, gender, country ];
	};

	const fetchCountries = async (countries) => {
		const response = await fetch(`https://restcountries.com/v3.1/alpha?codes=${countries}`);
		const country_names = await response.json();
		return country_names;
	};

	//-------------------------------------Promise.all-----------------------------------------//
	/* 
	const predictPerson = (name) => {
		document.getElementById('name-form').reset();
		if (name === '') {
			document.getElementsByClassName('predict-box')[0].style.display = 'none';
			return;
		}
		name = name.charAt(0).toUpperCase() + name.slice(1).toLowerCase();
		document.getElementById('name').textContent = ` ${name}`;
		let person = new Person();
		person.name = name;
		name.toLowerCase;
		Promise.all([
			fetch(`https://api.agify.io/?name=${name}`),
			fetch(`https://api.genderize.io?name=${name}`),
			fetch(`https://api.nationalize.io/?name=${name}`)
		])
			.then((responses) => {
				// Get a JSON object from each of the responses
				return Promise.all(responses.map((response) => response.json()));
			})
			.then((data) => {
				person.age = data[0].age;
				person.gender = data[1].gender;
				person.gender = person.gender.charAt(0).toUpperCase() + person.gender.slice(1);
				person.country = [];
				document.getElementById('gender').textContent = ` ${person.gender}`;
				document.getElementById('age').textContent = ` ${person.age}`;
				data[2].country.forEach((element) => person.country.push(element.country_id));
				document.getElementById('nationality').textContent = ` ${person.country}`;
				document.getElementsByClassName('predict-box')[0].style.display = 'flex';
			})
			.catch((error) => {
				// if there's an error, log it
				console.log(error);
			});
	};*/

	//-------------------------------------Chaining promises-----------------------------------------//

	/* 
	let predictPerson = (name) => {
		name.toLowerCase;
		document.getElementById('name').textContent = name.charAt(0).toUpperCase() + name.slice(1);
		let nation_url = 'https://api.nationalize.io/?name=' + name;
		let age_url = 'https://api.agify.io/?name=' + name;
		let gender_url = 'https://api.genderize.io?name=' + name;
		let country_url = 'https://restcountries.com/v3.1/alpha?codes=';

		let homosapien = new Person();
		homosapien.name = name.charAt(0).toUpperCase() + name.slice(1);
		let country_code = '';
		let country_name = '';

		fetch(age_url).then((response) => response.json()).then((data) => {
			homosapien.age = data.age;
			document.getElementById('age').textContent = data.age;
		});

		fetch(gender_url).then((response) => response.json()).then((data) => {
			homosapien.gender = data.gender;
			homosapien.gender = homosapien.gender.charAt(0).toUpperCase() + homosapien.gender.slice(1);
			homosapien.gender_probability = data.probability;
			document.getElementById('gender').textContent = ` ${homosapien.gender}`;
			country_code = country_code.replace(/,\s*$/, '');
		});

		fetch(nation_url)
			.then((response) => response.json())
			.then((data) => {
				homosapien.country = data.country;
			})
			.then(() => {
				homosapien.country.forEach((element, idx) => (country_code += `${element.country_id},`));
			})
			.then(() => fetch(`${country_url}${country_code}`))
			.then((response) => response.json())
			.then((data) => {
				data.forEach((element) => (country_name += `${element.name.common}, `));
				country_name = country_name.replace(/,\s*$/, '');
				document.getElementById('nationality').textContent = ` ${country_name}`;
			});

		setTimeout(() => {
			document.getElementsByClassName('predict-box')[0].style.display = 'flex';
			console.log(homosapien);
			console.log(country_name);
		}, 0);
	};*/
};
