import java.util.stream.*;

class Test {
    
    public void main(String[] args) {
        IntStream.range(1, 19)
            .mapToObj(String::valueOf)
            .forEach(System.out::println);
        IntStream.range(1, 2)
            .mapToObj(String::valueOf)
            .forEach(System.out::println);
        IntStream.range(1, 2)
            .mapToObj(x -> String.valueOf(x * 2))
            .forEach(System.out::println);
    }
}