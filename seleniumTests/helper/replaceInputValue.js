import { By } from 'selenium-webdriver';
export async function replaceInputValue(driver, name, value) {
	const input = await driver.findElement(By.name(name));
	await input.clear();
	await input.sendKeys(value);
}
