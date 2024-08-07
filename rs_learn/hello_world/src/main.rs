fn main() {
    // 变量的可变性
    let mut age:u8 = 17;
    age+=1;
    print!("age is:{}",age);

    // _
    let _name = "Tom";

    // 变量遮蔽
    let a1 = "a1";
    let a1 = 1;
    println!("value is: {}",a1);
}
