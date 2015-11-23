import java.util.*;


public class AttachmentIndex {

    public static void Main(string[] args) {
        var uri = new Uri("http://127.0.0.1:9200/");

        var connectionSettings = new ConnectionSettings(uri, defaultIndex:"my-application");
        var elasticClient = new ElasticClient(connectionSettings);
        // ref: http://nest.azurewebsites.net/nest/connecting.html


        elasticClient.CreateIndex("pdf-index", c => c.AddMapping<Document>(m => m.MapFromAttributes()));

        var attachment = new Attachment
        {
            Content = Convert.ToBase64String(File.ReadAllBytes("test.pdf")),
            ContentType = "application/pdf",
            Name = "test.pdf"
        };

        var doc = new Document()
        {
            Id = 1,
            Title = "test",
            File = attachment
        };

        elasticClient.Index(doc);
    }

    
}


