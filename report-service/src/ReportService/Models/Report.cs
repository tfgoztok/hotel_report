using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace ReportService.Models
{
    public class Report
    {
        // Unique identifier for the report
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string Id { get; set; } = null!;

        // Date when the report was requested
        public DateTime RequestDate { get; set; }
        
        // Current status of the report
        public string Status { get; set; } = null!;
        
        // Location associated with the report
        public string Location { get; set; } = null!;
        
        // Number of hotels included in the report
        public int HotelCount { get; set; }
        
        // Number of phone numbers included in the report
        public int PhoneNumberCount { get; set; }
    }
}
