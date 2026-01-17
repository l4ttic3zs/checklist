class FoodType {
  final int id;
  final String name;
  final String imagePath;

  FoodType({required this.id, required this.name, required this.imagePath});

  factory FoodType.fromJson(Map<String, dynamic> json) {
    return FoodType(
      id: json['id'],
      name: json['name'],
      imagePath: json['image_path'],
    );
  }
}