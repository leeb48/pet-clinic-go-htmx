import { By } from 'selenium-webdriver';
import { clickBtn } from '../helper/clickBtn.js';
import { getRandomRange } from './randomRange.js';

export async function goToRandomOwner(driver) {
	// go to owner list page

	await clickBtn(driver, 'owner-link');

	// go to first owner
	const tableBody = await driver.findElement(By.css('tbody'));
	const owners = await tableBody.findElements(By.css('tr'));
	const ownerIdx = getRandomRange(1, owners.length);
	const ownerBtn = await owners.at(ownerIdx).findElement(By.css('a'));

	const ownerDetailAction = driver.actions({ async: true });
	await ownerDetailAction.move({ origin: ownerBtn }).click().perform();
}
