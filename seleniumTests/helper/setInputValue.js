import { By } from 'selenium-webdriver';
export async function setInputValue(driver, name, value) {
	await driver.findElement(By.name(name)).sendKeys(value);
}
