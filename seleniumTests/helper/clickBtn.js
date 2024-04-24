import { By } from 'selenium-webdriver';
export async function clickBtn(driver, id) {
	const btn = await driver.findElement(By.id(id));
	const action = driver.actions({ async: true });
	await action.move({ origin: btn }).click().perform();
	await driver.manage().setTimeouts({ implicit: 300 });
}
