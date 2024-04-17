import assert from 'assert';
import { Browser, Builder, By, Capabilities } from 'selenium-webdriver';
import { addPet } from '../helper/addPet.js';
import { clickBtn } from '../helper/clickBtn.js';
import { goToRandomOwner } from '../helper/goToFirstOwner.js';

(async function removePet() {
	let driver;

	try {
		driver = await new Builder()
			.forBrowser(Browser.CHROME)
			.withCapabilities(
				Capabilities.chrome().set('acceptInsecureCerts', true)
			)
			.build();
		await driver.get('https://localhost:7771');

		await goToRandomOwner(driver);
		await clickBtn(driver, 'owner-edit-btn');
		await addPet(driver);
		await clickBtn(driver, 'owner-edit-submit');
		await driver.manage().setTimeouts({ implicit: 1000 });

		await clickBtn(driver, 'owner-edit-btn');
		await clickBtn(driver, 'pet-delete-btn');
		await clickBtn(driver, 'owner-edit-submit');
		await driver.manage().setTimeouts({ implicit: 1000 });

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
