import assert from 'assert';
import { Browser, Builder, By, Capabilities } from 'selenium-webdriver';
import {
	NumberDictionary,
	adjectives,
	animals,
	names,
	uniqueNamesGenerator,
} from 'unique-names-generator';
import { addPet } from '../helper/addPet.js';
import { clickBtn } from '../helper/clickBtn.js';
import { getRandomRange } from '../helper/randomRange.js';
import { goToRandomOwner } from '../helper/goToFirstOwner.js';
import { replaceInputValue } from '../helper/replaceInputValue.js';
import { setInputValue } from '../helper/setInputValue.js';

(async function editOwner() {
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

		await goToRandomOwner(driver);

		// go to edit page
		await clickBtn(driver, 'owner-edit-btn');

		// edit the owner
		const [firstName, lastName] = uniqueNamesGenerator({
			dictionaries: [names, names],
		}).split('_');

		await replaceInputValue(driver, 'firstName', firstName);
		await replaceInputValue(driver, 'lastName', lastName);

		const email = `${firstName.toLowerCase()}.${lastName.toLowerCase()}@test.com`;
		await replaceInputValue(driver, 'email', email);

		const phone = NumberDictionary.generate({
			min: 7021111111,
			max: 7029999999,
		})[0];
		await replaceInputValue(driver, 'phone', phone);

		const [month, day, year] = uniqueNamesGenerator({
			dictionaries: [
				NumberDictionary.generate({ min: 1, max: 12 }),
				NumberDictionary.generate({ min: 1, max: 31 }),
				NumberDictionary.generate({ min: 1970, max: 2015 }),
			],
		}).split('_');

		await replaceInputValue(driver, 'birthdate', `${month}/${day}/${year}`);

		const address = uniqueNamesGenerator({
			dictionaries: [
				NumberDictionary.generate({ min: 111, max: 9999 }),
				adjectives,
				animals,
			],
			separator: ' ',
		});
		await replaceInputValue(driver, 'address', address);

		const city = uniqueNamesGenerator({
			dictionaries: [names],
		});
		await replaceInputValue(driver, 'city', city);

		const stateInput = await driver.findElement(By.name('state'));
		const stateOptions = await stateInput.findElements(By.css('option'));
		const stateIdx = getRandomRange(1, stateOptions.length);
		await setInputValue(
			driver,
			'state',
			await stateOptions.at(stateIdx).getText()
		);

		// add pets
		await addPet(driver);

		await clickBtn(driver, 'owner-edit-submit');

		// check for valid response
		const alert = driver.findElement(By.className('alert alert-primary'));
		const response = await alert.getText();

		assert.equal(response, 'Edit Successful');
	} catch (e) {
		console.log(e);
	} finally {
		await driver.quit();
	}
})();
