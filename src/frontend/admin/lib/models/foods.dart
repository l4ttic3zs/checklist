import 'food_types.dart';

class Foods {
  final int id;
  final FoodType foodtype;
  final int count;

  Foods({required this.id, required this.foodtype, required this.count});

  factory Foods.fromJson(Map<String, dynamic> json) {
    return Foods(
      id: json['id'],
      foodtype: json['item_type'],
      count: json['count'],
    );
  }
}