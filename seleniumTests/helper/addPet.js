import { By } from 'selenium-webdriver';
import {
	NumberDictionary,
	names,
	uniqueNamesGenerator,
} from 'unique-names-generator';
import { getRandomRange } from './randomRange.js';

export async function addPet(driver) {
	const petName = uniqueNamesGenerator({
		dictionaries: [names],
	});

	await driver.findElement(By.id('pet-name-input')).sendKeys(petName);

	const petTypeInput = await driver.findElement(By.id('pet-type-input'));
	const options = await petTypeInput.findElements(By.css('option'));
	const randomSelect = getRandomRange(1, options.length);

	const selectedPetType = await options.at(randomSelect).getText();
	await driver.findElement(By.id('pet-type-input')).sendKeys(selectedPetType);

	const [month, day, year] = uniqueNamesGenerator({
		dictionaries: [
			NumberDictionary.generate({ min: 1, max: 12 }),
			NumberDictionary.generate({ min: 1, max: 31 }),
			NumberDictionary.generate({ min: 2000, max: 2020 }),
		],
	}).split('_');

	await driver
		.findElement(By.id('pet-birthdate-input'))
		.sendKeys(`${month}/${day}/${year}`);

	const addPetForm = driver.actions({ async: true });
	const addPetBtn = await driver.findElement(By.id('add-pet-btn'));
	await addPetForm.move({ origin: addPetBtn }).click().perform();
}
