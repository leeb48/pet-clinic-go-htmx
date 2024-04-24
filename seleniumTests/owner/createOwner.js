import assert from 'assert';
import { Browser, Builder, By, Capabilities } from 'selenium-webdriver';
import {
	NumberDictionary,
	adjectives,
	animals,
	names,
	uniqueNamesGenerator,
} from 'unique-names-generator';
import { clickBtn } from '../helper/clickBtn.js';
import { setInputValue } from '../helper/setInputValue.js';
import { addPet } from '../helper/addPet.js';

(async function createOwner() {
	let driver;

	try {
		driver = await new Builder()
			.forBrowser(Browser.CHROME)
			.withCapabilities(
				Capabilities.chrome().set('acceptInsecureCerts', true)
			)
			.build();
		await driver.get('https://localhost:7771');

		// go to owner list page
		await clickBtn(driver, 'owner-link');
		await driver.manage().setTimeouts({ implicit: 1000 });

		// go to create owner page
		await clickBtn(driver, 'owner-create-btn');

		// fill out the form
		const [firstName, lastName] = uniqueNamesGenerator({
			dictionaries: [names, names],
		}).split('_');

		await setInputValue(driver, 'firstName', firstName);
		await setInputValue(driver, 'lastName', lastName);

		const email = `${firstName.toLowerCase()}.${lastName.toLowerCase()}@test.com`;
		await setInputValue(driver, 'email', email);

		const phone = NumberDictionary.generate({
			min: 7021111111,
			max: 7029999999,
		})[0];
		setInputValue(driver, 'phone', phone);

		const [month, day, year] = uniqueNamesGenerator({
			dictionaries: [
				NumberDictionary.generate({ min: 1, max: 12 }),
				NumberDictionary.generate({ min: 1, max: 31 }),
				NumberDictionary.generate({ min: 1970, max: 2015 }),
			],
		}).split('_');
		await setInputValue(driver, 'birthdate', `${month}/${day}/${year}`);

		const address = uniqueNamesGenerator({
			dictionaries: [
				NumberDictionary.generate({ min: 111, max: 9999 }),
				adjectives,
				animals,
			],
			separator: ' ',
		});
		await setInputValue(driver, 'address', address);

		await setInputValue(driver, 'city', 'Las Vegas');
		await setInputValue(driver, 'state', 'Nevada');

		await addPet(driver);
		await addPet(driver);

		// submit the form
		await clickBtn(driver, 'owner-create-submit');

		// check for valid response
		const alert = await driver.findElement(
			By.className('alert alert-primary')
		);

		const response = await alert.getText();

		assert.equal(response, 'User created');
	} catch (e) {
		console.log(e);
	} finally {
		await driver.quit();
	}
})();
